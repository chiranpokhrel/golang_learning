package gorumspaxos

import (
	"sync"
	"time"

	pb "dat520/lab5/gorumspaxos/proto"
)

const (
	// promiseTimeout to wait for Promise messages
	promiseTimeout = 5 * time.Second
	// learnTimeout to wait for Learn messages
	learnTimeout = 5 * time.Second
	// requestWaitTime is a delay to wait for more requests if there are none
	requestWaitTime = 500 * time.Millisecond
)

// Proposer represents a proposer as defined by the Multi-Paxos algorithm.
type Proposer struct {
	mu                 sync.RWMutex
	id                 int               // replica's id.
	leader             int               // current Paxos leader.
	crnd               Round             // replica's current round; initially, this replica's id.
	adu                Slot              // all-decided-up-to is the highest consecutive slot that has been committed.
	nextSlot           Slot              // slot for the next request, initially 0.
	phaseOneDone       bool              // indicates if the phase1 is done, initially false.
	config             MultiPaxosConfig  // configuration used for multipaxos.
	nodeMap            map[string]uint32 // map of the address to the node id.
	acceptMsgQueue     []*pb.AcceptMsg   // queue of pending accept messages as part of prepare operation.
	clientRequestQueue []*pb.AcceptMsg   // queue of pending client requests.
}

// NewProposer returns a new Multi-Paxos proposer with the specified
// replica id, initial leader, and nodeMap.
func NewProposer(myID, leader int, nodeMap map[string]uint32) *Proposer {
	propIdx := myIndex(myID, nodeMap)
	return &Proposer{
		id:                 myID,
		leader:             leader,
		nodeMap:            nodeMap,
		crnd:               Round(propIdx),
		acceptMsgQueue:     make([]*pb.AcceptMsg, 0),
		clientRequestQueue: make([]*pb.AcceptMsg, 0),
	}
}

// newLeader updates the current leader and crnd.
func (p *Proposer) newLeader(leader int) {
	// TODO(student): complete
}

// isLeader returns true if this replica is the leader.
func (p *Proposer) isLeader() bool {
	// TODO(student): complete
	return false
}

// advanceAllDecidedUpTo increments the highest consecutive slot that has been committed.
func (p *Proposer) advanceAllDecidedUpTo() uint32 {
	// TODO(student): complete
	return p.adu
}

// runPhaseOne runs MultiPaxos phase one (prepare->promise).
//
// This method should only be called if this replica is the leader and phase one
// has not already completed.
//
// Steps:
//  1. Create a PrepareMsg with the current round and slot adu+1.
//  2. Send the PrepareMsg to all the replicas via the Prepare quorum call.
//  3. Process the combined promise message from the Prepare call.
//  4. For each accepted PValue in the promise message, prepare an AcceptMsg
//     and add it to the accept queue.
//  5. Advance the nextSlot to adu+1.
//  6. Set phaseOneDone to true.
func (p *Proposer) runPhaseOne() error {
	// TODO(student): complete
	return nil
}

// runMultiPaxos runs MultiPaxos phase one and two.
//
// This method should only be called if this replica is the leader and phase one
// has not already completed.
//
// If phase 1 has not been completed: Run phase 1
// Otherwise:
//
//	Check if there are pending requests in the clientRequestQueue or acceptMsgQueue
//	Otherwise: Wait for requestWaitTime then return
//	Call performAccept
//	Call performCommit with the returned learn message
func (p *Proposer) runMultiPaxos() {
	// TODO(student): complete
}

// nextAcceptMsg returns the next accept message to be sent, if any.
// If there are no pending accept messages or any client requests to process,
// it returns nil.
func (p *Proposer) nextAcceptMsg() (accept *pb.AcceptMsg) {
	// TODO(student): complete
	return accept
}

// Perform the accept quorum call on the replicas.
//
//  1. Check if any pending accept requests in the acceptReqQueue to process
//  2. Check if any pending client requests in the clientRequestQueue to process
//  3. Increment the nextSlot and prepare an accept message for the pending request,
//     using crnd and nextSlot.
//  4. Perform accept quorum call on the configuration and return the learnMsg.
func (p *Proposer) performAccept(accept *pb.AcceptMsg) (*pb.LearnMsg, error) {
	// TODO(student): complete
	return nil, nil
}

// Perform the commit operation using a multicast call.
func (p *Proposer) performCommit(learn *pb.LearnMsg) error {
	// TODO(student): complete
	_ = learn // TODO: remove this line when you implement the method
	return nil
}

// isPhaseOneDone return true if phase one is done.
func (p *Proposer) isPhaseOneDone() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.phaseOneDone
}

// setConfiguration set the configuration for the proposer, allowing it to
// communicate with the other replicas.
func (p *Proposer) setConfiguration(config MultiPaxosConfig) {
	p.mu.Lock()
	p.config = config
	p.mu.Unlock()
}

// AddRequestToQ adds the request to the clientRequestQueue.
func (p *Proposer) AddRequestToQ(request *pb.Value) {
	if p.isLeader() {
		p.Logf("Adding request to queue: %v", request)
		accept := &pb.AcceptMsg{Val: request}
		p.mu.Lock()
		p.clientRequestQueue = append(p.clientRequestQueue, accept)
		p.mu.Unlock()
	}
}
