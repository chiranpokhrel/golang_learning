package gorumspaxos

import (
	"context"
	"slices"

	pb "dat520/lab5/gorumspaxos/proto"

	"github.com/relab/gorums"
)

// myIndex returns the index of myID in the sorted list of IDs from the node map.
func myIndex(myID int, nodeMap map[string]uint32) int {
	// TODO(meling): When Go 1.23 is released, we can use the following code:
	// ids := slices.Sorted(maps.Values(nodeMap))
	ids := Values(nodeMap)
	slices.Sort(ids)
	return slices.Index(ids, uint32(myID))
}

func Keys[K comparable, V any](m map[K]V) []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func Values[K comparable, V any](m map[K]V) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

type (
	Round = int32
	Slot  = uint32
)

const (
	NoRound int32  = -1
	NoSlot  uint32 = 0
)

// Example values to be used in testing of acceptor and proposer.
var (
	valOne = &pb.Value{
		ClientID:      "1234",
		ClientSeq:     42,
		ClientCommand: "ls",
	}
	valTwo = &pb.Value{
		ClientID:      "5678",
		ClientSeq:     99,
		ClientCommand: "rm",
	}
	valThree = &pb.Value{
		ClientID:      "1369",
		ClientSeq:     4,
		ClientCommand: "mkdir",
	}
)

// MultiPaxosConfig defines the RPC calls on the configuration.
// This interface is used for mocking the configuration in unit tests.
type MultiPaxosConfig interface {
	Prepare(ctx context.Context, request *pb.PrepareMsg) (response *pb.PromiseMsg, err error)
	Accept(ctx context.Context, request *pb.AcceptMsg) (response *pb.LearnMsg, err error)
	Commit(ctx context.Context, request *pb.LearnMsg, opts ...gorums.CallOption)
	ClientHandle(ctx context.Context, request *pb.Value) (response *pb.Response, err error)
}

// Mock configuration used for testing.
// It stores the received value and allows you to specify the returned value for all RPC calls.
type MockConfiguration struct {
	pb.Configuration

	ErrOut error

	PrpIn  *pb.PrepareMsg
	PrmOut *pb.PromiseMsg

	AccIn  *pb.AcceptMsg
	LrnOut *pb.LearnMsg

	LrnIn  *pb.LearnMsg
	EmpOut *pb.Empty

	ValIn   *pb.Value
	RespOut *pb.Response
}

func (mc *MockConfiguration) Prepare(ctx context.Context, request *pb.PrepareMsg) (response *pb.PromiseMsg, err error) {
	mc.PrpIn = request
	return mc.PrmOut, mc.ErrOut
}

func (mc *MockConfiguration) Accept(ctx context.Context, request *pb.AcceptMsg) (response *pb.LearnMsg, err error) {
	mc.AccIn = request
	return mc.LrnOut, mc.ErrOut
}

func (mc *MockConfiguration) Commit(ctx context.Context, request *pb.LearnMsg, opts ...gorums.CallOption) {
	mc.LrnIn = request
}

func (mc *MockConfiguration) ClientHandle(ctx context.Context, request *pb.Value) (response *pb.Response, err error) {
	mc.ValIn = request
	return mc.RespOut, mc.ErrOut
}

type mockLD struct{}

func (mld *mockLD) Subscribe() <-chan int {
	ch := make(chan int)
	return ch
}

func (mld *mockLD) Leader() int {
	return 0
}
