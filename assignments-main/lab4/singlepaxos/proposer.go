package singlepaxos

// Proposer represents a proposer as defined by the single-decree Paxos algorithm.
type Proposer struct {
	crnd        Round
	clientValue Value
	// TODO(student): algorithm implementation
	// Add other needed fields
}

// NewProposer returns a new single-decree Paxos proposer.
// It takes the following arguments:
//
// id: The id of the node running this instance of a Paxos proposer.
//
// numNodes: The total number of Paxos nodes.
//
// The proposer's internal crnd field should initially be set to the value of
// its id.
func NewProposer(id int, numNodes int) *Proposer {
	// TODO(student): algorithm implementation
	return &Proposer{}
}

// handlePromise processes the promise according to the single-decree Paxos algorithm.
// It returns an accept message to send if the proposer has gathered a majority of promises.
// If an empty accept message is returned, the proposer should ignore the promise.
func (p *Proposer) handlePromise(promise Promise) Accept {
	// TODO(student): algorithm implementation
	return Accept{From: -1, Rnd: -2, Val: "FooBar"}
}

// increaseCrnd increases proposer p's crnd field by the total number
// of Paxos nodes.
func (p *Proposer) increaseCrnd() {
	// TODO(student): algorithm implementation
}

// TODO(student): Add any other unexported methods needed.
