package singlepaxos

// Acceptor represents an acceptor as defined by the single-decree Paxos algorithm.
type Acceptor struct { // TODO(student): algorithm implementation
	// Add needed fields
}

// NewAcceptor returns a new single-decree Paxos acceptor.
// It takes the following arguments:
//
// id: The id of the node running this instance of a Paxos acceptor.
func NewAcceptor(id int) *Acceptor {
	// TODO(student): algorithm implementation
	return &Acceptor{}
}

// handlePrepare processes the prepare according to the single-decree Paxos algorithm,
// returning a promise, or an empty promise if the prepare should be ignored.
func (a *Acceptor) handlePrepare(prepare Prepare) Promise {
	// TODO(student): algorithm implementation
	return Promise{To: -1, From: -1, Vrnd: -2, Vval: "FooBar"}
}

// handleAccept processes the accept according to the single-decree Paxos algorithm,
// returning a learn, or an empty learn if the accept should be ignored.
func (a *Acceptor) handleAccept(accept Accept) Learn {
	// TODO(student): algorithm implementation
	return Learn{From: -1, Rnd: -2, Val: "FooBar"}
}

// TODO(student): Add any other unexported methods needed.
