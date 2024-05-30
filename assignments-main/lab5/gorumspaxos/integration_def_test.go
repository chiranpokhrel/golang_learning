package gorumspaxos

import (
	"context"
	"fmt"
	"net"
	"slices"
	"sync"
	"testing"
	"time"

	pb "dat520/lab5/gorumspaxos/proto"

	"github.com/google/go-cmp/cmp"
	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	waitForReplicasToStart = 1000 * time.Millisecond
	waitForReplicasToStop  = 1000 * time.Millisecond
	waitForClientsToStop   = 1000 * time.Millisecond
	waitTimeForRequest     = 5 * time.Second
)

func startReplicas(t testing.TB, numServers int) (map[string]uint32, func(), func() int, []*PaxosReplica) {
	t.Helper()
	nodeMap := make(map[string]uint32)
	lisMap := make(map[string]net.Listener)
	for i := range numServers {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		addr := lis.Addr().String()
		nodeMap[addr] = uint32(i)
		lisMap[addr] = lis
	}
	replicas := make([]*PaxosReplica, 0)
	for addr, id := range nodeMap {
		replica := NewPaxosReplica(int(id), nodeMap)
		replicas = append(replicas, replica)
		go func() { replica.Serve(lisMap[addr]) }()
	}
	stopFn := func() {
		for _, replica := range replicas {
			replica.Stop()
		}
		t.Log("All replicas have been stopped")
	}
	crashLeader := func() int {
		for _, replica := range replicas {
			if replica.isLeader() {
				replica.Stop()
				t.Logf("Leader %d stopped", replica.id)
				return replica.leader
			}
		}
		return -1
	}
	time.Sleep(waitForReplicasToStart)
	return nodeMap, stopFn, crashLeader, replicas
}

func newConfiguration(nodeMap map[string]uint32) (*pb.Configuration, func(), error) {
	mgr := pb.NewManager(gorums.WithDialTimeout(5*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(), // block until connections are made
			grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
		),
	)
	qspec := NewPaxosQSpec(len(nodeMap))
	cfg, err := mgr.NewConfiguration(qspec, gorums.WithNodeMap(nodeMap))
	return cfg, func() {
		mgr.Close()
	}, err
}

func testFiveReplicas(t *testing.T, handler func()) {
	numReplicas := 5
	numClients := 3
	numRequests := 5

	nodeMap, teardown, _, _ := startReplicas(t, numReplicas)
	defer teardown()

	t.Log("All replicas have started")

	// start clients in separate goroutines
	var wg sync.WaitGroup
	wg.Add(numClients)
	errChan := make(chan error, numClients)
	for i := range numClients {
		go func(id int) {
			defer wg.Done()
			config, closeMgr, err := newConfiguration(nodeMap)
			if err != nil {
				errChan <- err
				return
			}
			defer closeMgr()

			for k := range numRequests {
				ctx, cancel := context.WithTimeout(context.Background(), waitTimeForRequest)
				defer cancel()
				req := &pb.Value{
					ClientID:      fmt.Sprint(id),
					ClientSeq:     uint32(k),
					ClientCommand: fmt.Sprint(k),
				}
				wantResp := &pb.Response{
					ClientID:      req.ClientID,
					ClientSeq:     req.ClientSeq,
					ClientCommand: req.ClientCommand,
				}
				gotResp, err := config.ClientHandle(ctx, req)
				if err != nil {
					errChan <- err
					return
				}
				if diff := cmp.Diff(wantResp, gotResp, protocmp.Transform()); diff != "" {
					errChan <- fmt.Errorf("%s mismatch (-want +got):\n%s", t.Name(), diff)
					return
				}
			}
		}(i)
	}
	// close the error channel once all goroutines are done
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		if err != nil {
			t.Error(err)
		}
	}
	if !t.Failed() {
		handler()
	}
}

func testLeaderFailure(t *testing.T, handler func()) {
	numReplicas := 5
	numClients := 3

	nodeMap, teardown, crash, replicas := startReplicas(t, numReplicas)

	// count number of correctly processed messages across all clients
	var lock sync.Mutex
	numMsgs := 0

	// start clients in separate goroutines
	var wg sync.WaitGroup
	wg.Add(numClients)
	errChan := make(chan error, numClients)
	errLogChan := make(chan string, 20)
	stopChan := make(chan bool, numClients)
	for i := range numClients {
		go func(id int) {
			defer wg.Done()
			config, closeMgr, err := newConfiguration(nodeMap)
			if err != nil {
				errChan <- err
				return
			}
			defer closeMgr()

			// continuously send requests to replicas
			for k := 0; true; k++ {
				select {
				case <-stopChan:
					return
				default:
					ctx, cancel := context.WithTimeout(context.Background(), waitTimeForRequest)
					defer cancel()
					req := &pb.Value{
						ClientID:      fmt.Sprint(id),
						ClientSeq:     uint32(k),
						ClientCommand: fmt.Sprint(k),
					}
					resp, err := config.ClientHandle(ctx, req)
					if err != nil {
						// errors that occur during leader change are expected; log and continue
						if len(errLogChan) < cap(errLogChan) {
							// only log the first few errors
							errLogChan <- fmt.Sprintf("ClientHandle(%v): err %v", req, err)
						}
						continue
					}
					if !req.Match(resp) {
						errChan <- fmt.Errorf("ClientHandle(%v) mismatch:\n%s", req, resp)
						return // stop this client goroutine if there is a mismatch
					}
					lock.Lock()
					numMsgs++
					lock.Unlock()
				}
			}
		}(i)
	}

	<-time.After(2 * time.Second)
	t.Log("Crashing leader replica")
	crashedReplica := crash()
	<-time.After(waitForReplicasToStop)

	lock.Lock()
	numMsgsAtCrash := numMsgs
	lock.Unlock()

	<-time.After(7 * time.Second)
	t.Log("Stopping clients")
	for range numClients {
		stopChan <- true
	}

	// close the error channel once all goroutines are done
	go func() {
		wg.Wait()
		close(errChan)
		close(errLogChan)
	}()
	// report client errors if any
	for err := range errChan {
		if err != nil {
			t.Error(err)
		}
	}
	t.Log("Logging client errors due to leader change; these are expected")
	for err := range errLogChan {
		if err != "" {
			t.Log(err)
		}
	}
	<-time.After(waitForClientsToStop)

	teardown()
	<-time.After(waitForReplicasToStop)

	lock.Lock()
	switch numMsgs {
	case 0:
		t.Errorf("Liveness failure: no messages were processed by the replicas")
	case numMsgsAtCrash:
		t.Errorf("Failed to perform leader change: no messages were processed after the leader crashed")
	}
	t.Logf("Total msgs / msgs before leader crash = %d / %d", numMsgs, numMsgsAtCrash)
	lock.Unlock()

	slices.SortFunc(replicas, func(ri, rj *PaxosReplica) int { return int(ri.id - rj.id) })

	checkAllCompleted(t, replicas, crashedReplica, numMsgsAtCrash, numMsgs)
	checkReplicaState(t, replicas)
	logRemainingResponses(t, replicas)
	if !t.Failed() {
		handler()
	}
}

func logRemainingResponses(t *testing.T, replicas []*PaxosReplica) {
	t.Helper()
	for _, replica := range replicas {
		remainingResponses := replica.remainingResponses()
		if remainingResponses > 0 {
			t.Logf("Replica[%d] has %d remaining responses to send", replica.id, remainingResponses)
		}
		if remainingResponses < 10 {
			for _, k := range replica.responseIDs() {
				t.Logf("Replica[%d] has remaining responses to send: %d", replica.id, k)
			}
		}
	}
}

func checkAllCompleted(t *testing.T, replicas []*PaxosReplica, crashedReplica int, numMsgsAtCrash int, numMsgs int) {
	t.Helper()
	for _, replica := range replicas {
		wantMsgs := want(replica.id, crashedReplica, numMsgsAtCrash, numMsgs)
		if len(replica.learntVal) != wantMsgs {
			t.Errorf("len(replica[%d].learntVal) = %d, want %d", replica.id, len(replica.learntVal), wantMsgs)
		}
	}
	t.Log("-----------------")
	for _, replica := range replicas {
		t.Logf("Replica[%d] has completed %d slots", replica.id, len(replica.learntVal))
	}
}

// want returns equal if i == j, otherwise different.
func want[T comparable, V any](i, j T, equal, different V) V {
	if i == j {
		return equal
	}
	return different
}

func checkReplicaState(t *testing.T, replicas []*PaxosReplica) {
	t.Helper()
	for _, r1 := range replicas {
		for _, r2 := range replicas {
			i, j := r1.id, r2.id
			if i < j {
				slot := sharedPrefix(r1.learntVal, r2.learntVal)
				if slot != NoSlot {
					t.Errorf("Replica[%d] != Replica[%d] at slot %d:\n\t%v\n\t%v", i, j, slot, r1.learntVal[slot], r2.learntVal[slot])
				}
			}
		}
	}
}

// sharedPrefix returns the first slot where the two maps differ.
func sharedPrefix(m1, m2 map[Slot]*pb.LearnMsg) Slot {
	shortestMap := m1
	if len(m2) < len(shortestMap) {
		shortestMap = m2
	}
	prefixSlots := Keys(shortestMap)
	slices.Sort(prefixSlots)
	for _, slot := range prefixSlots {
		if !m1[slot].Equal(m2[slot]) {
			return slot
		}
	}
	return NoSlot
}
