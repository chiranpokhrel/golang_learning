package gorumsfd

import (
	"slices"
	"sync"
	"testing"
	"time"

	pb "dat520/lab3/gorumsfd/proto"
	"dat520/lab3/leaderdetector"

	"github.com/google/go-cmp/cmp"
	"github.com/relab/gorums"
)

var crashNewLeaderTests = []struct {
	name         string
	nodeIDs      []int
	crashedNodes []int
	wantLeaders  []int
}{
	{name: "Crashed 3, want leader 2", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{3}, wantLeaders: []int{3, 2}},
	{name: "Crashed 2, want leader 3", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{2}, wantLeaders: []int{3}},
	{name: "Crashed 1, want leader 3", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{1}, wantLeaders: []int{3}},
	{name: "Crashed 0, want leader 3", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{0}, wantLeaders: []int{3}},
	{name: "Crashed {3, 2}, want leader 1", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{3, 2}, wantLeaders: []int{3, 2, 1}},
	{name: "Crashed {2, 3}, want leader 1", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{2, 3}, wantLeaders: []int{3, 1}},
	{name: "Crashed {1, 2, 3}, want leader 0", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{1, 2, 3}, wantLeaders: []int{3, 0}},
	{name: "Crashed {3, 2, 1}, want leader 0", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{3, 2, 1}, wantLeaders: []int{3, 2, 1, 0}},
	{name: "Crashed {3, 2, 1, 0}, want no leader", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{3, 2, 1, 0}, wantLeaders: []int{3, 2, 1, 0, -1}},
	{name: "Crashed {0, 1, 2, 3}, want no leader", nodeIDs: []int{0, 1, 2, 3}, crashedNodes: []int{0, 1, 2, 3}, wantLeaders: []int{3, -1}},
	{name: "Crashed {2, 5}, want leader 4", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{2, 5}, wantLeaders: []int{5, 4}},
	{name: "Crashed {4, 5}, want leader 3", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{4, 5}, wantLeaders: []int{5, 3}},
	{name: "Crashed {5, 4}, want leader 3", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{5, 4}, wantLeaders: []int{5, 4, 3}},
	{name: "Crashed {5, 2, 4}, want leader 3", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{5, 2, 4}, wantLeaders: []int{5, 4, 3}},
	{name: "Crashed {5, 4, 3}, want leader 2", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{5, 4, 3}, wantLeaders: []int{5, 4, 3, 2}},
	{name: "Crashed {5, 2, 3, 4}, want leader 1", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{5, 2, 3, 4}, wantLeaders: []int{5, 4, 1}},
	{name: "Crashed {5, 3, 2, 4}, want leader 1", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{5, 3, 2, 4}, wantLeaders: []int{5, 4, 1}},
	{name: "Crashed {5, 4, 3, 2}, want leader 1", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{5, 4, 3, 2}, wantLeaders: []int{5, 4, 3, 2, 1}},
	{name: "Crashed {5, 4, 3, 2, 1}, want leader 0", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{5, 4, 3, 2, 1}, wantLeaders: []int{5, 4, 3, 2, 1, 0}},
	{name: "Crashed {5, 4, 3, 2, 1, 0}, want no leader", nodeIDs: []int{0, 1, 2, 3, 4, 5}, crashedNodes: []int{5, 4, 3, 2, 1, 0}, wantLeaders: []int{5, 4, 3, 2, 1, 0, -1}},
}

func TestCrashNewLeader(t *testing.T) {
	const delta time.Duration = 5 * time.Millisecond

	for _, test := range crashNewLeaderTests {
		t.Run(test.name, func(t *testing.T) {
			fds := make(map[int]FailureDetector)

			// Mutex to protect the map of leader changes per node
			mu := sync.Mutex{}
			// Map of leader changes per node
			gotLeadersMap := make(map[int][]int)
			// Each node's done channel is used to stop waiting for leader changes.
			// This initialization must be done before starting the goroutines below to
			// avoid a data race.
			done := make(map[int]chan struct{})
			for _, id := range test.nodeIDs {
				done[id] = make(chan struct{})
				gotLeadersMap[id] = make([]int, 0)
			}

			wg := sync.WaitGroup{}
			wg.Add(len(test.nodeIDs))
			for _, id := range test.nodeIDs {
				ld := leaderdetector.NewMonLeaderDetector(test.nodeIDs)
				fds[id] = NewGorumsFailureDetector(uint32(id), ld, delta)
				go func(id int) {
					mu.Lock()
					gotLeadersMap[id] = append(gotLeadersMap[id], ld.Leader())
					mu.Unlock()
					leaderChan := ld.Subscribe()
					for {
						select {
						case gotLeader := <-leaderChan:
							mu.Lock()
							gotLeadersMap[id] = append(gotLeadersMap[id], gotLeader)
							mu.Unlock()

						case <-done[id]:
							wg.Done()
							return
						}
					}
				}(id)
			}

			// Start all failure detectors and generate heartbeats
			for _, fd := range fds {
				hbSender := func(hb *pb.HeartBeat) {
					for _, fd := range fds {
						fd.Heartbeat(gorums.ServerCtx{}, hb)
					}
				}
				fd.Start(hbSender)
			}
			<-time.After(2 * delta)

			// Crash selected nodes
			for _, id := range test.crashedNodes {
				fds[id].Stop()
				<-time.After(5 * delta)
			}

			// Stop all failure detectors
			for id, fd := range fds {
				fd.Stop()
				done[id] <- struct{}{}
			}
			wg.Wait() // for all to stop

			// Check the leader changes for each node
			for id, gotLeaders := range gotLeadersMap {
				// Determine if the node is in the crashed list.
				cIndex := slices.Index(test.crashedNodes, id)
				wantLeaders := test.wantLeaders

				// If the node crashed, adjust the expected leaders.
				if cIndex >= 0 {
					wIndex := slices.Index(test.wantLeaders, id)
					if wIndex >= 0 {
						// If the node was an expected leader, include leaders up to and including this node.
						wantLeaders = test.wantLeaders[:wIndex+1]
					} else {
						// If the node was not an expected leader, include the leaders that crashed before this node.
						wantLeaders = leadersUpToCrash(test.wantLeaders, test.crashedNodes)
					}
				}
				if diff := cmp.Diff(gotLeaders, wantLeaders); diff != "" {
					t.Errorf("Node %d: (-got +want)\n%s", id, diff)
				}
			}
		})
	}
}

// leadersUpToCrash returns a slice of leaders up to the first crash.
func leadersUpToCrash(leaders, crashed []int) []int {
	ret := []int{leaders[0]} // Always include the first leader.
	for _, leaderID := range leaders[1:] {
		if slices.Index(crashed, leaderID) != -1 {
			ret = append(ret, leaderID)
		}
	}
	return ret
}

type heartbeatEvent struct {
	desc          string
	heartbeats    []*pb.HeartBeat
	wantSuspected []int
	wantRestored  []int
	wantDelay     time.Duration
}

var heartbeatTests = []struct {
	nodes      []int
	delay      time.Duration
	heartbeats []heartbeatEvent
}{
	{
		nodes: []int{0, 1, 2, 3}, delay: 1, heartbeats: []heartbeatEvent{
			{desc: "Node 3 suspected", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 1}, {ID: 2}}, wantSuspected: []int{3}, wantRestored: []int{}, wantDelay: 1},
			{desc: "Node 3 restored", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}}, wantSuspected: []int{}, wantRestored: []int{3}, wantDelay: 2},
			{desc: "All restored ok", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}}, wantSuspected: []int{}, wantRestored: []int{}, wantDelay: 2},
			{desc: "Nodes {1, 2} suspected", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 3}}, wantSuspected: []int{1, 2}, wantRestored: []int{}, wantDelay: 2},
			{desc: "Node 2 restored", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 2}, {ID: 3}}, wantSuspected: []int{}, wantRestored: []int{2}, wantDelay: 3},
			{desc: "Node 1 restored", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}}, wantSuspected: []int{}, wantRestored: []int{1}, wantDelay: 4},
			{desc: "All restored ok", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}}, wantSuspected: []int{}, wantRestored: []int{}, wantDelay: 4},
			{desc: "Node 0 suspected", heartbeats: []*pb.HeartBeat{{ID: 1}, {ID: 2}, {ID: 3}}, wantSuspected: []int{0}, wantRestored: []int{}, wantDelay: 4},
			{desc: "Nodes {0, 1} suspected", heartbeats: []*pb.HeartBeat{{ID: 2}, {ID: 3}}, wantSuspected: []int{1}, wantRestored: []int{}, wantDelay: 4},
			{desc: "Nodes {0, 1, 2} suspected", heartbeats: []*pb.HeartBeat{{ID: 3}}, wantSuspected: []int{2}, wantRestored: []int{}, wantDelay: 4},
			{desc: "Nodes {0, 1, 2} suspected II", heartbeats: []*pb.HeartBeat{{ID: 3}}, wantSuspected: []int{}, wantRestored: []int{}, wantDelay: 4},
			{desc: "Nodes {0, 1, 2, 3} suspected", heartbeats: []*pb.HeartBeat{}, wantSuspected: []int{3}, wantRestored: []int{}, wantDelay: 4},
			{desc: "Nodes {1, 2} suspected, {0, 3} restored", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 3}}, wantSuspected: []int{}, wantRestored: []int{0, 3}, wantDelay: 5},
			{desc: "Nodes {1, 2} restored", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}}, wantSuspected: []int{}, wantRestored: []int{1, 2}, wantDelay: 6},
			{desc: "All restored ok", heartbeats: []*pb.HeartBeat{{ID: 0}, {ID: 1}, {ID: 2}, {ID: 3}}, wantSuspected: []int{}, wantRestored: []int{}, wantDelay: 6},
		},
	},
}

func TestHeartbeats(t *testing.T) {
	for _, test := range heartbeatTests {
		ld := &mockLD{
			nodes:     test.nodes,
			suspected: []int{},
			restore:   []int{},
		}
		fd := NewGorumsFailureDetector(uint32(0), ld, test.delay)
		for _, hbSeq := range test.heartbeats {
			for _, hb := range hbSeq.heartbeats {
				fd.Heartbeat(gorums.ServerCtx{}, hb)
			}
			fd.timeout()
			gotSuspected := ld.suspected
			if diff := cmp.Diff(gotSuspected, hbSeq.wantSuspected); diff != "" {
				t.Errorf("%s: Suspect (-got +want)\n%s", hbSeq.desc, diff)
			}
			gotRestored := ld.restore
			if diff := cmp.Diff(gotRestored, hbSeq.wantRestored); diff != "" {
				t.Errorf("%s: Restore (-got +want)\n%s", hbSeq.desc, diff)
			}
			gotDelay := fd.delay
			if diff := cmp.Diff(gotDelay, hbSeq.wantDelay); diff != "" {
				t.Errorf("%s: Delay (-got +want)\n%s", hbSeq.desc, diff)
			}
			ld.suspected = []int{}
			ld.restore = []int{}
		}
	}
}
