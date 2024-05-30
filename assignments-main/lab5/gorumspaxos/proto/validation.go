package proto

import (
	"crypto/sha256"
	"encoding/binary"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

type (
	Round = int32
	Slot  = uint32
)

const (
	NoRound int32  = -1
	NoSlot  uint32 = 0
)

// IsValid returns true if the promise message corresponds to the prepare message
// and the accepted pvalues are valid.
func (prepare *PrepareMsg) IsValid(promise *PromiseMsg) bool {
	switch {
	case prepare == nil && promise != nil:
		return false
	case prepare != nil && promise == nil:
		return false
	case prepare == nil && promise == nil:
		return true
	}
	for _, pval := range promise.Accepted {
		if !pval.IsValid(prepare) {
			// The acceptor should not have accepted the value in pval.
			// We ignore invalid promise messages.
			return false
		}
	}
	return prepare.Crnd == promise.Rnd
}

// IsValid returns true if pval is a valid response to the prepare message.
func (pval *PValue) IsValid(prepare *PrepareMsg) bool {
	switch {
	case pval == nil || prepare == nil:
		return false
	case prepare.Slot == NoSlot || prepare.Crnd == NoRound:
		return false // invalid input: prepare must be set
	case pval.Slot == NoSlot || pval.Vrnd == NoRound:
		return false // invalid pval; the acceptor should not have accepted a value without an assigned slot and round number
	case pval.Slot < prepare.Slot:
		return false // pval's slot is old
	case pval.Vval == nil:
		return false // pval's value is nil
	}
	// if pval.Vrnd > prepare.Crnd it means that the acceptor has already
	// accepted a value for a higher round than prepare.Crnd. This can only
	// happen if the acceptor has received another prepare with a higher round
	// than the one in prepare.Crnd. Hence, the acceptor is faulty because it
	// shouldn't have replied with a promise message. Hence, we only return
	// true if pval.Vrnd <= prepare.Crnd.
	return pval.Vrnd <= prepare.Crnd
}

// Match returns true if the learn message corresponds to the accept message.
func (accept *AcceptMsg) Match(learn *LearnMsg) bool {
	switch {
	case accept == nil && learn != nil:
		return false
	case accept != nil && learn == nil:
		return false
	case accept == nil && learn == nil:
		return true
	}
	return accept.Slot == learn.Slot &&
		accept.Rnd == learn.Rnd &&
		cmp.Equal(accept.Val, learn.Val, protocmp.Transform())
}

// Match returns true if the response corresponds to the request.
func (request *Value) Match(response *Response) bool {
	switch {
	case request == nil && response != nil:
		return false
	case request != nil && response == nil:
		return false
	case request == nil && response == nil:
		return true
	}
	return request.ClientID == response.ClientID &&
		request.ClientSeq == response.ClientSeq &&
		request.ClientCommand == response.ClientCommand
}

// Hash returns a hash of the request.
func (request *Value) Hash() uint64 {
	if request == nil {
		return 0
	}
	h := sha256.New()
	h.Write([]byte(request.ClientID))
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, request.ClientSeq)
	h.Write([]byte(b))
	h.Write([]byte(request.ClientCommand))
	bHash := h.Sum(nil)
	return binary.LittleEndian.Uint64(bHash)
}

// Equal returns true if the two learn messages are equal.
func (learn *LearnMsg) Equal(other *LearnMsg) bool {
	if learn == nil && other == nil {
		return true
	}
	if learn == nil || other == nil {
		return false
	}
	// We ignore the round number (Rnd) when comparing because the round number
	// may be different at different replicas.
	return learn.Slot == other.Slot && cmp.Equal(learn.Val, other.Val, protocmp.Transform())
}
