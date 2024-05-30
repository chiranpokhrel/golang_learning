package gorumspaxos

import (
	pb "dat520/lab5/gorumspaxos/proto"
)

// Acceptor represents an acceptor as defined by the Multi-Paxos algorithm.
type Acceptor struct {
	rnd         Round               // highest round the acceptor has promised in.
	accepted    map[Slot]*pb.PValue // map of accepted values for each slot.
	highestSeen Slot                // highest slot for which a prepare has been received.
}

// NewAcceptor returns a new Multi-Paxos acceptor
func NewAcceptor() *Acceptor {
	return &Acceptor{
		rnd:      NoRound,
		accepted: make(map[Slot]*pb.PValue),
	}
}

// handlePrepare processes the prepare according to the Multi-Paxos algorithm,
// returning a promise, or nil if the prepare should be ignored.
func (a *Acceptor) handlePrepare(prepare *pb.PrepareMsg) (prm *pb.PromiseMsg) {
	// TODO(student): Complete the handlePrepare
	return &pb.PromiseMsg{}
}

// handleAccept processes the accept according to the Multi-Paxos algorithm,
// returning a learn, or nil if the accept should be ignored.
func (a *Acceptor) handleAccept(accept *pb.AcceptMsg) (lrn *pb.LearnMsg) {
	// TODO(student): Complete the handleAccept
	return &pb.LearnMsg{}
}
