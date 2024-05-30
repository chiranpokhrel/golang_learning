package failuredetector

import "fmt"

type hbType bool

const (
	hbRequest = true
	hbReply   = false
)

// A Heartbeat is the basic message used by failure detectors to communicate
// with other nodes. A heartbeat message can be of type request or reply.
type Heartbeat struct {
	From int
	To   int
	Type hbType // true -> hbRequest, false -> hbReply
}

func (h Heartbeat) String() string {
	if h.Type {
		return fmt.Sprintf("Heartbeat request from %d to %d", h.From, h.To)
	}
	return fmt.Sprintf("Heartbeat reply from %d to %d", h.From, h.To)
}

// FailureDetector is the interface implemented by a failure detector.
type FailureDetector interface {
	Start()
	Stop()
	DeliverHeartbeat(hb Heartbeat)
}

// Suspecter is the interface that wraps the Suspect method. Suspect indicates
// that the node with identifier id should be considered suspected.
type Suspecter interface {
	Suspect(id int)
}

// Restorer is the interface that wraps the Restore method. Restore indicates
// that the node with identifier id should be considered restored.
type Restorer interface {
	Restore(id int)
}

// SuspectRestorer is the interface that groups the Suspect and Restore methods.
type SuspectRestorer interface {
	Suspecter
	Restorer
}
