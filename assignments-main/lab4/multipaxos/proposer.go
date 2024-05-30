package multipaxos

import (
	"dat520/lab3/leaderdetector"
)

// Proposer represents a proposer as defined by the Multi-Paxos algorithm.
type Proposer struct {
	id           int
	quorum       int
	n            int
	crnd         Round
	adu          Slot
	promises     []*Promise
	promiseCount int
	ld           leaderdetector.LeaderDetector
	leader       int
}

// NewProposer returns a new Multi-Paxos proposer. It takes the following
// arguments:
//
// id: The id of the node running this instance of a Paxos proposer.
//
// numNodes: The total number of Paxos nodes.
//
// adu: all-decided-up-to. The highest consecutive slot that has been decided.
// Should be set to -1 in the constructor, but for testing purposes we pass it
// to the constructor to be able to simulate the proposer's behavior when the
// proposer has already decided some slots.
//
// ld: A leader detector implementing the detector.LeaderDetector interface.
//
// The proposer's crnd field should initially be set to the value of its id.
func NewProposer(id, numNodes, adu int, ld leaderdetector.LeaderDetector) *Proposer {
	return &Proposer{
		id:       id,
		quorum:   (numNodes / 2) + 1,
		n:        numNodes,
		crnd:     Round(id),
		adu:      Slot(adu),
		promises: make([]*Promise, numNodes),
		ld:       ld,
		leader:   ld.Leader(),
	}
}

// handlePromise processes the promise according to the Multi-Paxos algorithm.
// It returns a slice of accept messages to send if the proposer has gathered
// a majority of promises, and should send accept messages. Accept messages
// whose Val field is the zero value is unconstrained and can be set to any value.
// If the slice is empty, the proposer is unconstrained and can send any value
// in accept messages. If nil is returned, the proposer should ignore the promise.
func (p *Proposer) handlePromise(prm Promise) []Accept {
	// TODO(student): algorithm implementation
	return []Accept{
		{From: -1, Slot: -1, Rnd: -2},
	}
}
