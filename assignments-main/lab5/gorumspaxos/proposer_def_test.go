package gorumspaxos

import (
	"errors"
	"slices"
	"testing"

	pb "dat520/lab5/gorumspaxos/proto"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func testPerformAccept(t *testing.T, handler func()) {
	for _, test := range performAcceptTests {
		proposer := newMockProposer(test.state)
		accept := proposer.nextAcceptMsg()
		gotLearn, err := proposer.performAccept(accept)
		if (err != nil) != test.wantError {
			t.Errorf("performAccept() error = %v, want %v", err, test.wantError)
		}
		if diff := cmp.Diff(gotLearn, test.wantLearn, protocmp.Transform()); diff != "" {
			t.Log(test.desc)
			t.Errorf("performAccept() mismatch (-want +got):\n%s", diff)
		}
		if !t.Failed() {
			handler()
		}
	}
}

func testPerformCommit(t *testing.T, handler func()) {
	for _, test := range performCommitTests {
		proposer := newMockProposer(test.state)
		err := proposer.performCommit(test.learn)
		if (err != nil) != test.wantError {
			t.Errorf("performCommit() error = %v, want %v", err, test.wantError)
		} else {
			handler()
		}
	}
}

func testRunPhaseOne(t *testing.T, handler func()) {
	// Ignore the slot of the accept message. Since they do not matter at this point.
	ignoreSlot := protocmp.IgnoreFields(&pb.AcceptMsg{}, "Slot")
	for _, test := range runPhaseOneTests {
		proposer := newMockProposer(test.state)
		err := proposer.runPhaseOne()
		mockConf := test.state.config
		acceptMsgs := slices.Clone(proposer.acceptMsgQueue)
		result := &expectedPhaseOne{
			Prepare:        mockConf.PrpIn,
			AcceptMsgQueue: acceptMsgs,
			PhaseOneDone:   proposer.phaseOneDone,
			Error:          err != nil,
		}
		if diff := cmp.Diff(result, test.expected, protocmp.Transform(), ignoreSlot); diff != "" {
			t.Log(test.desc)
			t.Errorf("runPhaseOne() mismatch (-want +got):\n%s", diff)
			if err != nil {
				t.Logf("Got error: %v", err)
			}
		} else {
			handler()
		}
	}
}

type proposerState struct {
	acceptMsgQueue []*pb.AcceptMsg
	clientMsgQueue []*pb.Value
	config         *MockConfiguration
	crnd           Round
	adu            Slot
	nextSlot       Slot
}

func newMockProposer(ps *proposerState) *Proposer {
	ld := &mockLD{}
	p := NewProposer(0, ld.Leader(), map[string]uint32{})
	p.setConfiguration(ps.config)
	p.acceptMsgQueue = slices.Clone(ps.acceptMsgQueue)
	p.clientRequestQueue = make([]*pb.AcceptMsg, len(ps.clientMsgQueue))
	for i, v := range ps.clientMsgQueue {
		p.clientRequestQueue[i] = &pb.AcceptMsg{Val: v}
	}
	p.crnd = ps.crnd
	p.adu = ps.adu
	p.nextSlot = ps.nextSlot
	return p
}

var performAcceptTests = []struct {
	desc      string
	state     *proposerState
	wantLearn *pb.LearnMsg
	wantError bool
}{
	{
		desc: "No messages in acceptMsgQueue or clientRequestQueue. Unable to perform accept. Expect nil learn.",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{},
			clientMsgQueue: []*pb.Value{},
			config:         &MockConfiguration{},
		},
		wantLearn: nil,
		wantError: false,
	},
	{
		desc: "One value in acceptMsgQueue. Use this to perform accept. Expect accept message with slot value incremented from proposer nextSlotID",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{{Rnd: 1, Slot: 0, Val: valOne}},
			clientMsgQueue: []*pb.Value{},
			config:         &MockConfiguration{LrnOut: &pb.LearnMsg{Rnd: 0, Val: valOne, Slot: 0}},
			crnd:           1,
		},
		wantLearn: &pb.LearnMsg{Slot: 0, Rnd: 0, Val: valOne},
		wantError: false,
	},
	{
		desc: "One value in clientRequestQueue. Use this to perform accept. Expect accept message with slot value incremented from proposer nextSlotID",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{},
			clientMsgQueue: []*pb.Value{valOne},
			config:         &MockConfiguration{LrnOut: &pb.LearnMsg{Rnd: 0, Val: valOne, Slot: 0}},
			crnd:           1,
		},
		wantLearn: &pb.LearnMsg{Slot: 0, Rnd: 0, Val: valOne},
		wantError: false,
	},
	{
		desc: "One value in clientRequestQueue and one value in acceptMsgQueue. Use message in acceptMsgQueue to perform accept. Expect accept message with slot value incremented from proposer nextSlotID",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{{Rnd: 1, Slot: 0, Val: valThree}},
			clientMsgQueue: []*pb.Value{valTwo},
			config:         &MockConfiguration{LrnOut: &pb.LearnMsg{Rnd: 1, Val: valThree, Slot: 1}},
			crnd:           1,
		},
		wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 1, Val: valThree},
		wantError: false,
	},
	{
		desc: "Accept message has old round id. Use the current round id of the proposer.",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{{Rnd: 1, Slot: 0, Val: valThree}},
			clientMsgQueue: []*pb.Value{valTwo},
			config:         &MockConfiguration{LrnOut: &pb.LearnMsg{Rnd: 10, Val: valThree, Slot: 6}},
			crnd:           10,
			nextSlot:       5,
		},
		wantLearn: &pb.LearnMsg{Slot: 6, Rnd: 10, Val: valThree},
		wantError: false,
	},
}

var performCommitTests = []struct {
	desc      string
	state     *proposerState
	learn     *pb.LearnMsg
	wantError bool
}{
	{
		desc: "Perform commit with provided learn message. Expect no error.",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{},
			clientMsgQueue: []*pb.Value{},
			config:         &MockConfiguration{},
		},
		learn:     &pb.LearnMsg{Slot: 6, Rnd: 10, Val: valThree},
		wantError: false,
	},
	{
		desc: "Perform commit without learn message. Expect error.",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{},
			clientMsgQueue: []*pb.Value{},
			config:         &MockConfiguration{},
		},
		learn:     nil,
		wantError: true,
	},
}

type expectedPhaseOne struct {
	Prepare        *pb.PrepareMsg
	AcceptMsgQueue []*pb.AcceptMsg
	PhaseOneDone   bool
	Error          bool
}

var runPhaseOneTests = []struct {
	desc     string
	state    *proposerState
	expected *expectedPhaseOne
}{
	{
		desc: "Send Prepare with incremented slot. Receive promise with no slots",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{},
			clientMsgQueue: []*pb.Value{},
			config:         &MockConfiguration{PrmOut: &pb.PromiseMsg{Rnd: 0, Accepted: []*pb.PValue{}}},
			crnd:           0,
			adu:            0,
		},
		expected: &expectedPhaseOne{
			Prepare:        &pb.PrepareMsg{Crnd: 0, Slot: 1},
			AcceptMsgQueue: []*pb.AcceptMsg{},
			PhaseOneDone:   true,
			Error:          false,
		},
	},
	{
		desc: "Send Prepare. Receive promise with a slot. Slot is properly added to acceptMsgQueue",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{},
			clientMsgQueue: []*pb.Value{},
			config:         &MockConfiguration{PrmOut: &pb.PromiseMsg{Rnd: 0, Accepted: []*pb.PValue{{Slot: 1, Vrnd: 0, Vval: valOne}}}},
			crnd:           0,
			adu:            1,
		},
		expected: &expectedPhaseOne{
			Prepare:        &pb.PrepareMsg{Crnd: 0, Slot: 2},
			AcceptMsgQueue: []*pb.AcceptMsg{{Rnd: 0, Val: valOne}},
			PhaseOneDone:   true,
			Error:          false,
		},
	},
	{
		desc: "Send Prepare. Receive promise with no slot.",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{},
			clientMsgQueue: []*pb.Value{},
			config:         &MockConfiguration{PrmOut: &pb.PromiseMsg{Rnd: 0, Accepted: []*pb.PValue{}}},
			crnd:           2,
			adu:            2,
		},
		expected: &expectedPhaseOne{
			Prepare:        &pb.PrepareMsg{Crnd: 2, Slot: 3},
			AcceptMsgQueue: []*pb.AcceptMsg{},
			PhaseOneDone:   true,
			Error:          false,
		},
	},
	{
		desc: "Send Prepare. Receive promise with three slots. All should be added to the acceptMsgQueue. The round value should be updated to the current round",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{},
			clientMsgQueue: []*pb.Value{},
			config: &MockConfiguration{
				PrmOut: &pb.PromiseMsg{
					Rnd: 2,
					Accepted: []*pb.PValue{
						{Slot: 1, Vrnd: 1, Vval: valThree},
						{Slot: 2, Vrnd: 1, Vval: valOne},
						{Slot: 4, Vrnd: 1, Vval: valTwo},
					},
				},
			},
			crnd: 3,
			adu:  2,
		},
		expected: &expectedPhaseOne{
			Prepare: &pb.PrepareMsg{Crnd: 3, Slot: 3},
			AcceptMsgQueue: []*pb.AcceptMsg{
				{Slot: 1, Rnd: 3, Val: valThree},
				{Slot: 2, Rnd: 3, Val: valOne},
				{Slot: 4, Rnd: 3, Val: valTwo},
			},
			PhaseOneDone: true,
			Error:        false,
		},
	},
	{
		desc: "Send Prepare. Error returned",
		state: &proposerState{
			acceptMsgQueue: []*pb.AcceptMsg{},
			clientMsgQueue: []*pb.Value{},
			config:         &MockConfiguration{PrmOut: nil, ErrOut: errors.New("test error")},
			crnd:           3,
			adu:            2,
		},
		expected: &expectedPhaseOne{
			Prepare:        &pb.PrepareMsg{Crnd: 3, Slot: 3},
			AcceptMsgQueue: []*pb.AcceptMsg{},
			PhaseOneDone:   false,
			Error:          true,
		},
	},
}
