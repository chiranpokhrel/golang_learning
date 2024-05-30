package gorumspaxos

import (
	"testing"

	pb "dat520/lab5/gorumspaxos/proto"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func testPaxosQSpec(t *testing.T, handler func()) {
	tests := []struct {
		size       int
		wantQuorum int
	}{
		{1, 1},
		{2, 1},
		{3, 2},
		{5, 3},
		{7, 4},
		{9, 5},
	}
	for _, test := range tests {
		qspecs := NewPaxosQSpec(test.size)
		gotQuorum := qspecs.quorum
		if gotQuorum != test.wantQuorum {
			t.Errorf("NewPaxosQSpec(%d) = %v, want %v", test.size, gotQuorum, test.wantQuorum)
		} else {
			handler()
		}
	}
}

func testPrepareQF(t *testing.T, handler func()) {
	for _, test := range prepareQFTests {
		if ok := t.Run(test.desc, func(t *testing.T) {
			qspecs := NewPaxosQSpec(test.configSize)
			replies := sliceToReplyMap(test.promises)
			gotPromise, gotQuorum := qspecs.PrepareQF(test.prepare, replies)
			if diff := cmp.Diff(gotPromise, test.wantPromise, protocmp.Transform()); diff != "" {
				t.Errorf("%s (-got +want)\n%s", test.desc, diff)
			}
			if gotQuorum != test.wantQuorum {
				t.Errorf("PrepareQF(%v, %v) = %v, %v; want %v, %v", test.prepare, replies, gotPromise, gotQuorum, test.wantPromise, test.wantQuorum)
			}
		}); ok {
			handler()
		}
	}
}

func testAcceptQF(t *testing.T, handler func()) {
	for _, test := range acceptQFTests {
		if ok := t.Run(test.desc, func(t *testing.T) {
			qspecs := NewPaxosQSpec(test.configSize)
			replies := sliceToReplyMap(test.learns)
			gotLearnMsg, gotQuorum := qspecs.AcceptQF(test.accept, replies)
			if diff := cmp.Diff(gotLearnMsg, test.wantLearn, protocmp.Transform()); diff != "" {
				t.Errorf("%s (-got +want)\n%s", test.desc, diff)
			}
			if gotQuorum != test.wantQuorum {
				t.Errorf("AcceptQF(%v, %v) = %v, %v; want %v, %v", test.accept, replies, gotLearnMsg, gotQuorum, test.wantLearn, test.wantQuorum)
			}
		}); ok {
			handler()
		}
	}
}

func testClientHandleQF(t *testing.T, handler func()) {
	for _, test := range clientHandleQFTests {
		if ok := t.Run(test.desc, func(t *testing.T) {
			qspecs := NewPaxosQSpec(test.configSize)
			replies := sliceToReplyMap(test.responses)
			gotResponse, gotQuorum := qspecs.ClientHandleQF(test.value, replies)
			if diff := cmp.Diff(gotResponse, test.wantResponse, protocmp.Transform()); diff != "" {
				t.Errorf("%s (-got +want)\n%s", test.desc, diff)
			}
			if gotQuorum != test.wantQuorum {
				t.Errorf("ClientHandleQF(%v, %v) = %v, %v; want %v, %v", test.value, replies, gotResponse, gotQuorum, test.wantResponse, test.wantQuorum)
			}
		}); ok {
			handler()
		}
	}
}

// sliceToReplyMap converts a slice of replies to a map of replies.
func sliceToReplyMap[V any](s []V) map[uint32]V {
	m := make(map[uint32]V)
	for i, v := range s {
		m[uint32(i)] = v
	}
	return m
}

var (
	noVal    = &pb.Value{}
	noopVal  = &pb.Value{IsNoop: true}
	oneVal   = &pb.Value{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"}
	twoVal   = &pb.Value{ClientID: "2", ClientSeq: 42, ClientCommand: "rm"}
	threeVal = &pb.Value{ClientID: "3", ClientSeq: 99, ClientCommand: "ps"}
	oneResp  = &pb.Response{ClientID: "1", ClientSeq: 1, ClientCommand: "rm"}
)

var prepareQFTests = []struct {
	desc        string
	configSize  int
	prepare     *pb.PrepareMsg   // input prepare message
	promises    []*pb.PromiseMsg // replies to the prepare message
	wantPromise *pb.PromiseMsg   // expected promise message from the quorum function
	wantQuorum  bool             // expected quorum result from the quorum function
}{
	{
		desc:       "Single promise response, no quorum",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 2, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 1}}},
		},
		wantPromise: nil,
		wantQuorum:  false,
	},
	{
		desc:       "Two valid promises, quorum, slot 2 already accepted in round 1",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 2, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 1, Vval: oneVal}}},
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 1, Vval: oneVal}}},
		},
		wantPromise: &pb.PromiseMsg{Rnd: 6},
		wantQuorum:  true,
	},
	{
		desc:       "Two valid promises, quorum, accepted slot 3 in round 1",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 2, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 3, Vrnd: 1, Vval: oneVal}}},
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 3, Vrnd: 1, Vval: oneVal}}},
		},
		wantPromise: &pb.PromiseMsg{Rnd: 6, Accepted: []*pb.PValue{{Slot: 3, Vrnd: 1, Vval: oneVal}}},
		wantQuorum:  true,
	},
	{
		desc:       "Two valid promises, quorum, accepted slot 3 in round 5",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 2, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 3, Vrnd: 5, Vval: oneVal}}},
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 3, Vrnd: 5, Vval: oneVal}}},
		},
		wantPromise: &pb.PromiseMsg{Rnd: 6, Accepted: []*pb.PValue{{Slot: 3, Vrnd: 5, Vval: oneVal}}},
		wantQuorum:  true,
	},
	{
		desc:       "Two promises, one invalid, no quorum, no promise",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 2, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 1}}},
			{Rnd: 3, Accepted: []*pb.PValue{{Slot: 2}}}, // invalid: old promise
		},
		wantPromise: nil,
		wantQuorum:  false,
	},
	{
		desc:       "Two valid promises, quorum, accepted slot 2 in round 5",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 1, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 5, Vval: oneVal}}},
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 5, Vval: oneVal}}},
		},
		wantPromise: &pb.PromiseMsg{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 5, Vval: oneVal}}},
		wantQuorum:  true,
	},
	{
		desc:       "Two valid promises, quorum, accepted slots {2,4} in round 5, missing slot 3",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 1, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 5, Vval: oneVal}, {Slot: 4, Vrnd: 5, Vval: twoVal}}},
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 1, Vrnd: 5, Vval: oneVal}, {Slot: 4, Vrnd: 5, Vval: twoVal}}},
		},
		wantPromise: &pb.PromiseMsg{
			Rnd: 6,
			Accepted: []*pb.PValue{
				{Slot: 2, Vrnd: 5, Vval: oneVal},
				{Slot: 3, Vrnd: 6, Vval: noopVal}, // fill the gap with a noop
				{Slot: 4, Vrnd: 5, Vval: twoVal},
			},
		},
		wantQuorum: true,
	},
	{
		desc:       "Two valid promises, quorum, accepted slots {2,3,4,6} in rounds 3 and 4, missing slot 5",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 1, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 3, Vval: oneVal}, {Slot: 4, Vrnd: 3, Vval: twoVal}, {Slot: 6, Vrnd: 3, Vval: threeVal}}},
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 1, Vrnd: 4, Vval: oneVal}, {Slot: 3, Vrnd: 4, Vval: twoVal}}},
		},
		wantPromise: &pb.PromiseMsg{
			Rnd: 6,
			Accepted: []*pb.PValue{
				{Slot: 2, Vrnd: 3, Vval: oneVal},
				{Slot: 3, Vrnd: 4, Vval: twoVal},
				{Slot: 4, Vrnd: 3, Vval: twoVal},
				{Slot: 5, Vrnd: 6, Vval: noopVal}, // fill the gap with a noop
				{Slot: 6, Vrnd: 3, Vval: threeVal},
			},
		},
		wantQuorum: true,
	},
	{
		desc:       "Two valid promises, quorum, accepted slots {2,3,6} in rounds 3 and 4, missing slot {4,5}",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 1, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 3, Vval: twoVal}, {Slot: 3, Vrnd: 4, Vval: twoVal}, {Slot: 6, Vrnd: 4, Vval: threeVal}}},
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 4, Vval: oneVal}, {Slot: 3, Vrnd: 3, Vval: oneVal}}},
		},
		wantPromise: &pb.PromiseMsg{
			Rnd: 6,
			Accepted: []*pb.PValue{
				{Slot: 2, Vrnd: 4, Vval: oneVal},
				{Slot: 3, Vrnd: 4, Vval: twoVal},
				{Slot: 4, Vrnd: 6, Vval: noopVal}, // fill the gap with a noop
				{Slot: 5, Vrnd: 6, Vval: noopVal}, // fill the gap with a noop
				{Slot: 6, Vrnd: 4, Vval: threeVal},
			},
		},
		wantQuorum: true,
	},
	{
		desc:       "Three promises, one invalid promise, quorum, accepted slots {2,3,6} in rounds 3 and 4, missing slot {4,5}",
		configSize: 3,
		prepare:    &pb.PrepareMsg{Slot: 1, Crnd: 6},
		promises: []*pb.PromiseMsg{
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 3, Vval: oneVal}, {Slot: 3, Vrnd: 4, Vval: twoVal}, {Slot: 6, Vrnd: 4, Vval: threeVal}}},
			{Rnd: 0, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 6, Vval: oneVal}, {Slot: 8, Vrnd: 4, Vval: twoVal}, {Slot: 7, Vrnd: 4, Vval: twoVal}}}, // invalid promise
			{Rnd: 6, Accepted: []*pb.PValue{{Slot: 2, Vrnd: 4, Vval: oneVal}, {Slot: 3, Vrnd: 3, Vval: twoVal}}},
		},
		wantPromise: &pb.PromiseMsg{
			Rnd: 6,
			Accepted: []*pb.PValue{
				{Slot: 2, Vrnd: 4, Vval: oneVal},
				{Slot: 3, Vrnd: 4, Vval: twoVal},
				{Slot: 4, Vrnd: 6, Vval: noopVal},
				{Slot: 5, Vrnd: 6, Vval: noopVal},
				{Slot: 6, Vrnd: 4, Vval: threeVal},
			},
		},
		wantQuorum: true,
	},
}

var acceptQFTests = []struct {
	desc       string
	configSize int
	accept     *pb.AcceptMsg  // input accept message
	learns     []*pb.LearnMsg // replies to the accept message
	wantLearn  *pb.LearnMsg   // expected learn message from the quorum function
	wantQuorum bool           // expected quorum result from the quorum function
}{
	{
		desc:       "One valid learn, no quorum",
		configSize: 3,
		accept:     &pb.AcceptMsg{Rnd: 3, Slot: 2, Val: oneVal},
		learns: []*pb.LearnMsg{
			{Rnd: 3, Slot: 2, Val: oneVal},
		},
		wantLearn:  nil,
		wantQuorum: false,
	},
	{
		desc:       "Two valid learns, quorum",
		configSize: 3,
		accept:     &pb.AcceptMsg{Rnd: 3, Slot: 2, Val: oneVal},
		learns: []*pb.LearnMsg{
			{Rnd: 3, Slot: 2, Val: oneVal},
			{Rnd: 3, Slot: 2, Val: oneVal},
		},
		wantLearn:  &pb.LearnMsg{Rnd: 3, Slot: 2, Val: oneVal},
		wantQuorum: true,
	},
	{
		desc:       "Three valid learns, quorum",
		configSize: 3,
		accept:     &pb.AcceptMsg{Rnd: 3, Slot: 2, Val: oneVal},
		learns: []*pb.LearnMsg{
			{Rnd: 3, Slot: 2, Val: oneVal},
			{Rnd: 3, Slot: 2, Val: oneVal},
			{Rnd: 3, Slot: 2, Val: oneVal},
		},
		wantLearn:  &pb.LearnMsg{Rnd: 3, Slot: 2, Val: oneVal},
		wantQuorum: true,
	},
	{
		desc:       "Two learns, one invalid slot, no quorum",
		configSize: 3,
		accept:     &pb.AcceptMsg{Rnd: 3, Slot: 2, Val: oneVal},
		learns: []*pb.LearnMsg{
			{Rnd: 3, Slot: 1, Val: oneVal}, // invalid learn: wrong slot
			{Rnd: 3, Slot: 2, Val: oneVal},
		},
		wantLearn:  nil,
		wantQuorum: false,
	},
	{
		desc:       "Two learns, one invalid round, no quorum",
		configSize: 3,
		accept:     &pb.AcceptMsg{Rnd: 3, Slot: 2, Val: oneVal},
		learns: []*pb.LearnMsg{
			{Rnd: 2, Slot: 2, Val: oneVal}, // invalid learn: wrong round
			{Rnd: 3, Slot: 2, Val: oneVal},
		},
		wantLearn:  nil,
		wantQuorum: false,
	},
	{
		desc:       "Two learns, one invalid value, no quorum",
		configSize: 3,
		accept:     &pb.AcceptMsg{Rnd: 3, Slot: 2, Val: oneVal},
		learns: []*pb.LearnMsg{
			{Rnd: 3, Slot: 2, Val: twoVal}, // invalid learn: wrong value
			{Rnd: 3, Slot: 2, Val: oneVal},
		},
		wantLearn:  nil,
		wantQuorum: false,
	},
	{
		desc:       "Three learns, one invalid slot, quorum",
		configSize: 3,
		accept:     &pb.AcceptMsg{Rnd: 3, Slot: 2, Val: oneVal},
		learns: []*pb.LearnMsg{
			{Rnd: 3, Slot: 1, Val: oneVal}, // invalid learn: wrong slot
			{Rnd: 3, Slot: 2, Val: oneVal},
			{Rnd: 3, Slot: 2, Val: oneVal},
		},
		wantLearn:  &pb.LearnMsg{Rnd: 3, Slot: 2, Val: oneVal},
		wantQuorum: true,
	},
	{
		desc:       "Three learns, one invalid round, quorum",
		configSize: 3,
		accept:     &pb.AcceptMsg{Rnd: 3, Slot: 2, Val: oneVal},
		learns: []*pb.LearnMsg{
			{Rnd: 2, Slot: 2, Val: oneVal}, // invalid learn: wrong round
			{Rnd: 3, Slot: 2, Val: oneVal},
			{Rnd: 3, Slot: 2, Val: oneVal},
		},
		wantLearn:  &pb.LearnMsg{Rnd: 3, Slot: 2, Val: oneVal},
		wantQuorum: true,
	},
	{
		desc:       "Three learns, one invalid value, quorum",
		configSize: 3,
		accept:     &pb.AcceptMsg{Rnd: 3, Slot: 2, Val: oneVal},
		learns: []*pb.LearnMsg{
			{Rnd: 3, Slot: 2, Val: twoVal}, // invalid learn: wrong value
			{Rnd: 3, Slot: 2, Val: oneVal},
			{Rnd: 3, Slot: 2, Val: oneVal},
		},
		wantLearn:  &pb.LearnMsg{Rnd: 3, Slot: 2, Val: oneVal},
		wantQuorum: true,
	},
}

var clientHandleQFTests = []struct {
	desc         string
	configSize   int
	value        *pb.Value      // input value message
	responses    []*pb.Response // replies to the value message
	wantResponse *pb.Response   // expected response message from the quorum function
	wantQuorum   bool           // expected quorum result from the quorum function
}{
	{
		desc:       "No quorum replies",
		configSize: 5,
		value:      oneVal,
		responses: []*pb.Response{
			oneResp,
			oneResp,
		},
		wantResponse: nil,
		wantQuorum:   false,
	},
	{
		desc:       "Quorum size (3) replies, all valid",
		configSize: 5,
		value:      oneVal,
		responses: []*pb.Response{
			oneResp,
			oneResp,
			oneResp,
		},
		wantResponse: oneResp,
		wantQuorum:   true,
	},
	{
		desc:       "Quorum size (3) replies, two valid, one with invalid ClientSeq, no quorum",
		configSize: 5,
		value:      oneVal,
		responses: []*pb.Response{
			oneResp,
			oneResp,
			{ClientID: "1", ClientSeq: 2, ClientCommand: "rm"},
		},
		wantResponse: nil,
		wantQuorum:   false,
	},
	{
		desc:       "Quorum size (3) replies, two valid, one with invalid ClientCommand, no quorum",
		configSize: 5,
		value:      oneVal,
		responses: []*pb.Response{
			oneResp,
			oneResp,
			{ClientID: "1", ClientSeq: 1, ClientCommand: "1rm"},
		},
		wantResponse: nil,
		wantQuorum:   false,
	},
	{
		desc:       "Quorum size (3) replies, two valid, one with invalid ClientID, no quorum",
		configSize: 5,
		value:      oneVal,
		responses: []*pb.Response{
			oneResp,
			oneResp,
			{ClientID: "2", ClientSeq: 1, ClientCommand: "rm"},
		},
		wantResponse: nil,
		wantQuorum:   false,
	},
	{
		desc:       "System size (3) replies, two valid, one with invalid ClientSeq, valid quorum",
		configSize: 3,
		value:      oneVal,
		responses: []*pb.Response{
			oneResp,
			oneResp,
			{ClientID: "1", ClientSeq: 2, ClientCommand: "rm"},
		},
		wantResponse: oneResp,
		wantQuorum:   true,
	},
	{
		desc:       "System size (3) replies, two valid, one with invalid ClientCommand, valid quorum",
		configSize: 3,
		value:      oneVal,
		responses: []*pb.Response{
			oneResp,
			oneResp,
			{ClientID: "1", ClientSeq: 1, ClientCommand: "1rm"},
		},
		wantResponse: oneResp,
		wantQuorum:   true,
	},
	{
		desc:       "System size (3) replies, two valid, one with invalid ClientId, valid quorum",
		configSize: 3,
		value:      oneVal,
		responses: []*pb.Response{
			oneResp,
			oneResp,
			{ClientID: "2", ClientSeq: 1, ClientCommand: "rm"},
		},
		wantResponse: oneResp,
		wantQuorum:   true,
	},
}
