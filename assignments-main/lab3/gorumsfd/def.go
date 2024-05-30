package gorumsfd

import (
	pb "dat520/lab3/gorumsfd/proto"
)

// FailureDetector is the interface implemented by a failure detector.
type FailureDetector interface {
	pb.FailureDetector
	Start(func(*pb.HeartBeat))
	Stop()
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

// SuspectRestorer is the interface that groups the Suspect and Restore
// methods.
type SuspectRestorer interface {
	Suspecter
	Restorer
	NodeIDs() []uint32
}

type mockLD struct {
	nodes     []int
	suspected []int
	restore   []int
}

func (ld *mockLD) Suspect(id int) {
	ld.suspected = append(ld.suspected, id)
}

func (ld *mockLD) Restore(id int) {
	ld.restore = append(ld.restore, id)
}

func (ld *mockLD) NodeIDs() []uint32 {
	nodeIDs := make([]uint32, len(ld.nodes))
	for i := range ld.nodes {
		nodeIDs[i] = uint32(ld.nodes[i])
	}
	return nodeIDs
}
