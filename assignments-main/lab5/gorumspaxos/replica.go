package gorumspaxos

import (
	"context"
	"dat520/lab3/gorumsfd"
	"dat520/lab3/leaderdetector"
	"errors"
	"net"
	"sync"
	"time"

	fd "dat520/lab3/gorumsfd/proto"
	pb "dat520/lab5/gorumspaxos/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// responseTimeout is the duration to wait for a response before cancelling
	responseTimeout = 1 * time.Second
	// managerDialTimeout is the default timeout for dialing a manager
	managerDialTimeout = 5 * time.Second
	// delta is the failure detector's default timeout value
	delta = 1 * time.Second
)

// PaxosReplica is the structure composing the Proposer and Acceptor.
type PaxosReplica struct {
	pb.MultiPaxos
	mu sync.Mutex
	*Acceptor
	*Proposer
	leaderDetector  leaderdetector.LeaderDetector
	failureDetector gorumsfd.FailureDetector
	fdManager       *fd.Manager             // gorums failure detector manager (from generated code)
	paxosManager    *pb.Manager             // gorums paxos manager (from generated code)
	id              int                     // id is the id of the node
	srv             *gorums.Server          // the gorums.Server that the replica is registered to
	stop            chan struct{}           // channel for stopping the replica's run loop.
	learntVal       map[uint32]*pb.LearnMsg // Stores all received learn messages
	stopped         bool
}

// NewPaxosReplica returns a new Paxos replica with a nodeMap configuration.
func NewPaxosReplica(myID int, nodeMap map[string]uint32) *PaxosReplica {
	nodeIds := make([]int, 0)
	for _, id := range nodeMap {
		nodeIds = append(nodeIds, int(id))
	}
	ld := leaderdetector.NewMonLeaderDetector(nodeIds)
	failureDetector := gorumsfd.NewGorumsFailureDetector(uint32(myID), ld, delta)

	opts := []gorums.ManagerOption{
		gorums.WithDialTimeout(managerDialTimeout),
		gorums.WithGrpcDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
		),
	}
	r := &PaxosReplica{
		Acceptor:        NewAcceptor(),
		Proposer:        NewProposer(myID, ld.Leader(), nodeMap),
		leaderDetector:  ld,
		failureDetector: failureDetector,
		fdManager:       fd.NewManager(opts...),
		paxosManager:    pb.NewManager(opts...),
		id:              myID,
		srv:             gorums.NewServer(),
		stop:            make(chan struct{}),
		learntVal:       make(map[uint32]*pb.LearnMsg),
	}
	fd.RegisterFailureDetectorServer(r.srv, r.failureDetector)
	pb.RegisterMultiPaxosServer(r.srv, r)
	r.run()
	return r
}

func newTestReplicaLeader() *PaxosReplica {
	myID := 0
	replica := &PaxosReplica{
		Proposer:  NewProposer(myID, myID, map[string]uint32{"0": 0}),
		id:        myID,
		learntVal: make(map[uint32]*pb.LearnMsg),
	}
	replica.Proposer.phaseOneDone = true
	replica.adu = 0
	return replica
}

// Stops the failure detector, replica, and gorums server.
func (r *PaxosReplica) Stop() {
	if r.stopped {
		return
	}
	r.stopped = true
	r.failureDetector.Stop()
	r.stop <- struct{}{} // stop the replica's run loop
	r.fdManager.Close()
	r.paxosManager.Close()
	r.srv.Stop()
}

// Serve starts the server and blocks until the server is stopped.
func (r *PaxosReplica) Serve(lis net.Listener) {
	if err := r.srv.Serve(lis); err != nil {
		r.Logf("Failed to serve: %v", err)
	}
}

// ServerStart starts the replica
// 1. Invokes the start function of the proposer
// 2. Create a new gorums server and store it to to the replica
// 3. Register MultiPaxos server
// 4. Start failure detector
// 5. Call Serve on gorums server
func (replica *PaxosReplica) ServerStart(lis net.Listener) {
	// TODO(student) Implement the function
}

// run starts the replica's run loop.
// It subscribes to the leader detector's trust messages and signals the proposer when a new leader is detected.
// It also starts the failure detector, which is necessary to get leader detections.
func (r *PaxosReplica) run() {
	trustMsgs := r.leaderDetector.Subscribe()
	go func() {
		// TODO(student) create Paxos configuration and set the proposer's configuration
		// TODO(student) implement Paxos replica run loop
		_ = trustMsgs // TODO: remove this line when you start using trustMsgs
		<-r.stop      // TODO: remove this line when you implement the method
	}()

	go func() {
		cfg, err := r.fdManager.NewConfiguration(gorums.WithNodeMap(r.nodeMap))
		if err != nil {
			r.Logf("Failed to create configuration for failure detector: %v", err)
			return
		}
		hbSender := func(hb *fd.HeartBeat) {
			cfg.Heartbeat(context.Background(), &fd.HeartBeat{ID: uint32(r.id)})
		}
		r.failureDetector.Start(hbSender)
	}()
}

// Prepare handles the prepare quorum calls from the proposer by passing the received messages to its acceptor.
// It receives prepare massages and pass them to handlePrepare method of acceptor.
// It returns promise messages back to the proposer by its acceptor.
func (r *PaxosReplica) Prepare(ctx gorums.ServerCtx, prepare *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	r.Logf("Acceptor: Prepare(%v) received", prepare)
	return r.handlePrepare(prepare), nil
}

// Accept handles the accept quorum calls from the proposer by passing the received messages to its acceptor.
// It receives Accept massages and pass them to handleAccept method of acceptor.
// It returns learn massages back to the proposer by its acceptor
func (r *PaxosReplica) Accept(ctx gorums.ServerCtx, accept *pb.AcceptMsg) (*pb.LearnMsg, error) {
	r.Logf("Acceptor: Accept(%v) received", accept)
	return r.handleAccept(accept), nil
}

// Commit is invoked by the proposer as part of the commit phase of the MultiPaxos algorithm.
// It receives a learn massage representing the proposer's decided value, meaning that the
// request can be executed by the replica. (In this lab you don't need to execute the request,
// just deliver the response to the client.)
//
// Be aware that the received learn message may not be for the next slot in the sequence.
// If the received slot is less than the next slot, the message should be ignored.
// If the received slot is greater than the next slot, the message should be buffered.
// If the received slot is equal to the next slot, the message should be delivered.
//
// This method is also responsible for communicating the decided value to the ClientHandle
// method, which is responsible for returning the response to the client.
func (r *PaxosReplica) Commit(ctx gorums.ServerCtx, learn *pb.LearnMsg) {
	r.Logf("Replica: Commit(%v) received", learn)
	r.mu.Lock()
	adu := r.adu + 1
	if _, ok := r.learntVal[learn.Slot]; !ok {
		r.learntVal[learn.Slot] = learn
	}
	r.mu.Unlock()
	_ = adu // TODO: remove this line when you implement the method
	// TODO(student) complete
}

// ClientHandle is invoked by the client to send a request to the replicas via a quorum call and get a response.
// A response is only sent back to the client when the request has been committed by the MultiPaxos replicas.
// This method will receive requests from multiple clients and must return the response to the correct client.
// If the request is not committed within a certain time, the method may return an error.
//
// Since the method is called by multiple clients, it is essential to return the matching reply to the client.
// Consider a client that sends a request M1, once M1 has been decided, the response to M1 should be returned
// to the client. However, while waiting for M1 to get committed, M2 may be proposed and committed by the replicas.
// Thus, M2 should not be returned to the client that sent M1.
func (r *PaxosReplica) ClientHandle(ctx gorums.ServerCtx, req *pb.Value) (rsp *pb.Response, err error) {
	r.AddRequestToQ(req)
	// TODO(student) complete
	return nil, errors.New("unable to get the response")
}

// remainingResponses returns the number of responses that are still pending.
func (r *PaxosReplica) remainingResponses() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	// TODO(student) complete
	return 0
}

// responseIDs returns the IDs of the responses that are still pending.
func (r *PaxosReplica) responseIDs() []uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	ids := make([]uint64, 0)
	// TODO(student) complete
	return ids
}
