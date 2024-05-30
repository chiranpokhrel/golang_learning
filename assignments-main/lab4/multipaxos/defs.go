package multipaxos

import "fmt"

// Slot represents an identifier for a Multi-Paxos consensus instance.
type Slot int

// Round represents a Multi-Paxos round number.
type Round int

// NoRound is a constant that represents no specific round. It should be used
// as the value for the Vrnd field in Promise messages to indicate that an
// acceptor has not voted in any previous round.
const NoRound Round = -1

// Value represents a value that can be chosen using the Multi-Paxos algorithm and
// has the following fields:
//
// ClientID: Unique identifier for the client that sent the command.
//
// ClientSeq: Client local sequence number.
//
// Command: The state machine command to be agreed upon and executed.
type Value struct {
	ClientID  string
	ClientSeq int
	Command   string
}

var noopValue = Value{}

// String returns a string representation of value v.
func (v Value) String() string {
	if v == noopValue {
		return "No-op value"
	}
	return fmt.Sprintf("Value{ClientID: %s, ClientSeq: %d, Command: %s}", v.ClientID, v.ClientSeq, v.Command)
}

// Prepare represents a Multi-Paxos prepare message.
type Prepare struct {
	From int
	Slot Slot
	Crnd Round
}

// String returns a string representation of prepare p.
func (p Prepare) String() string {
	return fmt.Sprintf("Prepare{From: %d, Slot: %d, Crnd: %d}", p.From, p.Slot, p.Crnd)
}

// Promise represents a Multi-Paxos promise message.
// The Accepted field is a set of PValues that have been accepted
// (by the acceptor that created the Promise) in a given slot.
type Promise struct {
	To, From int
	Rnd      Round
	Accepted []PValue
}

// String returns a string representation of promise p.
func (p Promise) String() string {
	if p.Accepted == nil {
		return fmt.Sprintf("Promise{To: %d, From: %d, Rnd: %d, No accepted values reported (nil slice)}", p.To, p.From, p.Rnd)
	}
	if len(p.Accepted) == 0 {
		return fmt.Sprintf("Promise{To: %d, From: %d, Rnd: %d, No accepted values reported (empty slice)}", p.To, p.From, p.Rnd)
	}
	return fmt.Sprintf("Promise{To: %d, From: %d, Rnd: %d, Accepted: %v}", p.To, p.From, p.Rnd, p.Accepted)
}

// Accept represents a Multi-Paxos Paxos accept message.
type Accept struct {
	From int
	Slot Slot
	Rnd  Round
	Val  Value
}

// String returns a string representation of accept a.
func (a Accept) String() string {
	return fmt.Sprintf("Accept{From: %d, Slot: %d, Rnd: %d, Val: %v}", a.From, a.Slot, a.Rnd, a.Val)
}

// Learn represents a Multi-Paxos learn message.
type Learn struct {
	From int
	Slot Slot
	Rnd  Round
	Val  Value
}

// String returns a string representation of learn l.
func (l Learn) String() string {
	return fmt.Sprintf("Learn{From: %d, Slot: %d, Rnd: %d, Val: %v}", l.From, l.Slot, l.Rnd, l.Val)
}

// PValue is a triple consisting of a round number, a slot number,
// and a value; the value is typically used to represent a command.
//
// A PValue is created when an acceptor votes for a value in a round and slot.
type PValue struct {
	Slot Slot
	Vrnd Round
	Vval Value
}
