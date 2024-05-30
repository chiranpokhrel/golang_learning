package singlepaxos

// Learner represents a learner as defined by the single-decree Paxos algorithm.
type Learner struct { // TODO(student): algorithm implementation
	// Add needed fields
	// Tip: you need to keep the decided values by the Paxos nodes somewhere
}

// NewLearner returns a new single-decree Paxos learner. It takes the
// following arguments:
//
// id: The id of the node running this instance of a Paxos learner.
//
// numNodes: The total number of Paxos nodes.
func NewLearner(id int, numNodes int) *Learner {
	// TODO(student): algorithm implementation
	return &Learner{}
}

// handleLearn processes the learn according to the single-decree Paxos algorithm,
// returning a value if the learn results in the learner emitting a decided value.
// Otherwise, it returns an empty value.
func (l *Learner) handleLearn(learn Learn) Value {
	// TODO(student): algorithm implementation
	return "FooBar"
}

// TODO(student): Add any other unexported methods needed.
