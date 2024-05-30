package main

import (
	"context"
	"testing"
	"time"

	"dat520/lab3/gorumsfd"
	pb "dat520/lab3/gorumsfd/proto"
	"dat520/lab3/leaderdetector"

	"github.com/relab/gorums"
)

func TestNewReplica(t *testing.T) {
	const delta time.Duration = 500 * time.Millisecond

	nodeIDs := []int{0, 1, 2, 3}
	replicas := make([]*Replica, len(nodeIDs))
	for _, id := range nodeIDs {
		ld := leaderdetector.NewMonLeaderDetector(nodeIDs)
		fd := gorumsfd.NewGorumsFailureDetector(uint32(id), ld, delta)
		r, err := NewReplica("", fd)
		if err != nil {
			t.Errorf("NewReplica() error = %v, want nil", err)
		}
		if r == nil {
			t.Error("NewReplica() = nil, want non-nil")
		}
		replicas[id] = r
		go func() {
			if err := r.Serve(); err != nil {
				t.Fatalf("Failed to start replica: %v", err)
			}
		}()
	}

	<-time.After(2 * delta)
	nodes := replicasToNodes(replicas)
	for _, r := range replicas {
		r.Start(nodes)
	}
	<-time.After(2 * delta)
	replicas[0].Stop()

	for i := 0; i < 3; i++ {
		replicas[0].Configuration().Heartbeat(context.Background(), &pb.HeartBeat{ID: 10})
	}

	<-time.After(5 * delta)
	for _, r := range replicas {
		r.Stop()
	}
}

func replicasToNodes(replicas []*Replica) gorums.NodeListOption {
	nodeMap := make(map[string]uint32)
	for i, r := range replicas {
		nodeMap[r.Addr()] = uint32(i)
	}
	return gorums.WithNodeMap(nodeMap)
}
