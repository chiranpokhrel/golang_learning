package multipaxos

// Learner represents a learner as defined by the Multi-Paxos algorithm.
type Learner struct {
	// TODO(student): algorithm implementation
	// Add needed fields
}

// NewLearner returns a new Multi-Paxos learner. It takes the following
// arguments:
//
// numNodes: The total number of Paxos nodes.
func NewLearner(numNodes int) *Learner {
	// TODO(student): algorithm implementation
	return &Learner{}
}

// handleLearn processes the learn according to the Multi-Paxos algorithm,
// returning the decided value for the slot, if a quorum of learns have been
// collected; otherwise, it returns an empty value and 0.
func (l *Learner) handleLearn(learn Learn) (Value, Slot) {
	// TODO(student): algorithm implementation
	return Value{ClientID: "-1", ClientSeq: -1, Command: "-1"}, -1
}

// TODO(student): Add any other unexported methods needed.
