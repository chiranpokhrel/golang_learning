package failuredetector

import "time"

// EvtFailureDetector represents a Eventually Perfect Failure Detector as
// described at page 53 in:
// Christian Cachin, Rachid Guerraoui, and Luís Rodrigues: "Introduction to
// Reliable and Secure Distributed Programming" Springer, 2nd edition, 2011.
type EvtFailureDetector struct {
	id        int          // the id of this node
	nodeIDs   []int        // node ids for every node in cluster
	alive     map[int]bool // map of node ids considered alive
	suspected map[int]bool // map of node ids considered suspected

	sr SuspectRestorer // Provided SuspectRestorer implementation

	delay         time.Duration // the current delay for the timeout procedure
	delta         time.Duration // the delta value to be used when increasing delay
	timeoutSignal *time.Ticker  // the timeout procedure ticker

	hbSend chan Heartbeat // channel for sending outgoing heartbeat messages
	hbRecv chan Heartbeat // channel for receiving incoming heartbeat messages
	stop   chan struct{}  // channel for signaling a stop request to the main run loop

	testHook func() // DO NOT REMOVE THIS LINE. A no-op when not testing.
}

// NewEvtFailureDetector returns a new Eventual Failure Detector. It takes the
// following arguments:
//
// id: The id of the node running this instance of the failure detector.
//
// nodeIDs: A list of ids for every node in the cluster (including the node
// running this instance of the failure detector).
//
// sr: A leader detector implementing the SuspectRestorer interface.
//
// delta: The initial value for the timeout interval. Also the value to be used
// when increasing delay.
func NewEvtFailureDetector(id int, nodeIDs []int, sr SuspectRestorer, delta time.Duration) *EvtFailureDetector {
	suspected := make(map[int]bool)
	alive := make(map[int]bool)

	// TODO(student): perform any initialization necessary

	return &EvtFailureDetector{
		id:        id,
		nodeIDs:   nodeIDs,
		alive:     alive,
		suspected: suspected,

		sr: sr,

		delay: delta,
		delta: delta,

		hbSend: make(chan Heartbeat, 16),
		hbRecv: make(chan Heartbeat, 8),
		stop:   make(chan struct{}),

		testHook: func() {}, // DO NOT REMOVE THIS LINE. A no-op when not testing.
	}
}

// Start starts e's main run loop as a separate goroutine. The main run loop
// handles incoming heartbeat requests and responses. The loop also trigger e's
// timeout procedure at an interval corresponding to e's internal delay
// duration variable.
func (e *EvtFailureDetector) Start() {
	e.timeoutSignal = time.NewTicker(e.delay)
	go func(timeout <-chan time.Time) {
		for {
			e.testHook() // DO NOT REMOVE THIS LINE. A no-op when not testing.
			select {
			case <-e.hbRecv:
				// TODO(student): Handle incoming heartbeat
			case <-timeout:
				e.timeout()
			case <-e.stop:
				return
			}
		}
	}(e.timeoutSignal.C)
}

// Stop stops the failure detector's main run loop.
//
// Calling Stop multiple times without a Start call in between them will panic.
func (e *EvtFailureDetector) Stop() {
	e.stop <- struct{}{}
}

// Deliver delivers heartbeat to the failure detector.
//
// After stopping the failure detector, Deliver will panic if called.
func (e *EvtFailureDetector) Deliver(heartbeat Heartbeat) {
	e.hbRecv <- heartbeat
}

// Heartbeats returns a receive-only channel on which the failure detector sends
// outgoing heartbeats.
//
//   - Heartbeat replies are sent in response to incoming heartbeat requests.
//   - Heartbeat requests are sent periodically to all other nodes in the cluster
//     as part of the timeout procedure.
//
// The channel is closed when the failure detector is stopped, after which no
// more heartbeats will be sent and trying to read from the channel will panic.
func (e *EvtFailureDetector) Heartbeats() <-chan Heartbeat {
	return e.hbSend
}

// timeout updates the failure detector's internal state and sends out
// heartbeat requests to all other nodes in the cluster. If the intersection
// between the set of alive nodes and the set of suspected nodes is not empty,
// the delay is increased by delta and the timeout procedure is restarted.
//
// After stopping the failure detector, timeout will panic if called.
func (e *EvtFailureDetector) timeout() {
	// TODO(student): Implement timeout procedure
}

// TODO(student): Add other unexported functions or methods if needed.
