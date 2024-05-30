package gorumsfd

import (
	"time"

	pb "dat520/lab3/gorumsfd/proto"

	"github.com/relab/gorums"
)

// GorumsFailureDetector is a variant of the Eventually Perfect Failure Detector.
type GorumsFailureDetector struct {
	myID      uint32          // the id of this node
	nodeIDs   []uint32        // list of node ids
	alive     map[uint32]bool // map of node ids considered alive
	suspected map[uint32]bool // map of node ids considered suspected
	sr        SuspectRestorer // Provided SuspectRestorer implementation
	delay     time.Duration   // the current delay for the timeout procedure
	delta     time.Duration   // the delta value to be used when increasing delay
	stop      chan struct{}   // channel for signaling a stop request to the main run loop
	// TODO(student) add more fields if needed
}

// NewGorumsFailureDetector returns a new Eventual Failure Detector. It takes the
// following arguments:
// myID: The id of the node running this instance of the failure detector.
// sr: A leader detector implementing the SuspectRestorer interface.
// delta: The initial timeout delay, and value to be used when increasing delay.
func NewGorumsFailureDetector(myID uint32, sr SuspectRestorer, delta time.Duration) *GorumsFailureDetector {
	return &GorumsFailureDetector{
		myID:      myID,
		nodeIDs:   sr.NodeIDs(),
		alive:     make(map[uint32]bool),
		suspected: make(map[uint32]bool),
		sr:        sr,
		delay:     delta,
		delta:     delta,
		stop:      make(chan struct{}),
	}
}

// Start starts the failure detector's run loop in a separate goroutine.
// This function should perform the following functionalities:
//  1. Periodically send heartbeats to all nodes in the configuration;
//     the provided hbSender function should be used for this purpose.
//  2. Periodically update the status of the nodes to the SuspectRestorer.
//  3. Wait to receive the signal to stop.
func (e *GorumsFailureDetector) Start(hbSender func(*pb.HeartBeat)) {
	go func() {
		// TODO(student) complete Start method
	}()
}

// Stop stops the failure detector.
// If the failure detector has already stopped, we return immediately.
func (e *GorumsFailureDetector) Stop() {
	// TODO(student) complete Stop method
}

// timeout updates the status of the nodes to the SuspectRestorer.
// If a previously suspected node is now alive then increase the e.delay by e.delta.
// All unreachable nodes are reported and all previously reported, now alive nodes
// are restored by calling the Suspect and Restore functions of the SuspectRestorer.
//
// The algorithm is described at page 53 in:
// Christian Cachin, Rachid Guerraoui, and Luís Rodrigues: "Introduction to
// Reliable and Secure Distributed Programming" Springer, 2nd edition, 2011.
func (e *GorumsFailureDetector) timeout() {
	// TODO(student) complete timeout method
}

// Heartbeat is a multicast call invoked on all nodes in the configuration.
func (e *GorumsFailureDetector) Heartbeat(ctx gorums.ServerCtx, in *pb.HeartBeat) {
	// TODO(student) complete Heartbeat method
}
