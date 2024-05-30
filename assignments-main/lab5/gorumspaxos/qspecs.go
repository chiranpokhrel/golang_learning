package gorumspaxos

import (
	pb "dat520/lab5/gorumspaxos/proto"
)

// PaxosQSpec is a quorum specification object for Paxos.
// It only holds the quorum size.
type PaxosQSpec struct {
	quorum int
}

// NewPaxosQSpec returns a quorum specification object for Paxos
// for the given configuration size n.
func NewPaxosQSpec(n int) PaxosQSpec {
	return PaxosQSpec{}
}

// PrepareQF is the quorum function to process the replies from the Prepare quorum call.
// This is where the Proposer handle PromiseMsgs returned by the Acceptors, and any
// Accepted values in the promise replies should be combined into the returned PromiseMsg
// in the increasing slot order with any gaps filled with no-ops. Invalid promise messages
// should be ignored. For example, quorum function should only process promise replies that
// belong to the current round, as indicated by the prepare message. The quorum function
// returns true if a quorum of valid promises was found, and the combined PromiseMsg.
// Nil and false is returned if no quorum of valid promises was found.
func (qs PaxosQSpec) PrepareQF(prepare *pb.PrepareMsg, replies map[uint32]*pb.PromiseMsg) (*pb.PromiseMsg, bool) {
	// TODO(student) complete the PrepareQF
	return nil, false
}

// AcceptQF is the quorum function to process the replies from the Accept quorum call.
// This is where the Proposer handle LearnMsgs to determine if a value has been decided
// by the Acceptors. The quorum function returns true if a value has been decided, and
// the corresponding LearnMsg holds the slot, round number and value that was decided.
// Nil and false is returned if no value was decided.
func (qs PaxosQSpec) AcceptQF(accept *pb.AcceptMsg, replies map[uint32]*pb.LearnMsg) (*pb.LearnMsg, bool) {
	// TODO (student) complete the AcceptQF
	return nil, false
}

// ClientHandleQF is the quorum function to process the replies from the ClientHandle quorum call.
// This is where the Client handle the replies from the replicas. The quorum function should
// validate the replies against the request, and only valid replies should be considered.
// The quorum function returns true if a quorum of the replicas replied with the same response,
// and a single response is returned. Nil and false is returned if no quorum of valid replies was found.
func (qs PaxosQSpec) ClientHandleQF(request *pb.Value, replies map[uint32]*pb.Response) (*pb.Response, bool) {
	// TODO (student) complete the ClientHandleQF
	return nil, false
}
