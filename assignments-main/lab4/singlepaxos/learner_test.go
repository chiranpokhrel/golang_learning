package singlepaxos

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSinglePaxosLearner(t *testing.T) {
	for _, test := range learnerTests {
		t.Run(test.name, func(t *testing.T) {
			for _, event := range test.learnerEvents {
				learn, wantValue := event.req, event.wantResp
				gotValue := test.learner.handleLearn(learn)
				if diff := cmp.Diff(wantValue, gotValue); diff != "" {
					t.Log(test.desc)
					t.Errorf("handleLearn() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

var learnerTests = []struct {
	name          string
	desc          string
	learner       *Learner
	learnerEvents []learnValue
}{
	{
		name:    "NoQuorum",
		desc:    "no quorum -> no output",
		learner: NewLearner(0, 3),
		learnerEvents: []learnValue{
			{
				req:      Learn{From: 1, Rnd: 1, Val: "Lamport"},
				wantResp: ZeroValue,
			},
		},
	},
	{
		name:    "Quorum",
		desc:    "two learns, 3 nodes, same round and value, unique senders = quorum -> report output and value",
		learner: NewLearner(0, 3),
		learnerEvents: []learnValue{
			{
				req:      Learn{From: 1, Rnd: 1, Val: "Lamport"},
				wantResp: ZeroValue,
			},
			{
				req:      Learn{From: 2, Rnd: 1, Val: "Lamport"},
				wantResp: "Lamport",
			},
		},
	},
	{
		name:    "NoQuorumSameSender",
		desc:    "two learns, 3 nodes, same round and value, same sender = no quorum -> no output",
		learner: NewLearner(0, 3),
		learnerEvents: []learnValue{
			{
				req:      Learn{From: 1, Rnd: 1, Val: "Lamport"},
				wantResp: ZeroValue,
			},
			{
				req:      Learn{From: 1, Rnd: 1, Val: "Lamport"},
				wantResp: ZeroValue,
			},
		},
	},
	{
		name:    "NoQuorumDifferentRounds",
		desc:    "two learns, 3 nodes, different rounds, unique senders = no quorum -> no output",
		learner: NewLearner(0, 3),
		learnerEvents: []learnValue{
			{
				req:      Learn{From: 1, Rnd: 1, Val: "Lamport"},
				wantResp: ZeroValue,
			},
			{
				req:      Learn{From: 2, Rnd: 2, Val: "Lamport"},
				wantResp: ZeroValue,
			},
		},
	},
	{
		name:    "IgnoreLowerRound",
		desc:    "two learns, 3 nodes, second learn should be ignored due to lower round -> no output",
		learner: NewLearner(0, 3),
		learnerEvents: []learnValue{
			{
				req:      Learn{From: 2, Rnd: 2, Val: "Lamport"},
				wantResp: ZeroValue,
			},
			{
				req:      Learn{From: 1, Rnd: 1, Val: "Lamport"},
				wantResp: ZeroValue,
			},
		},
	},
	{
		name:    "IgnoreDifferentValue",
		desc:    "two learns, 3 nodes, second learn should be ignored due to different value in same round -> no output",
		learner: NewLearner(0, 3),
		learnerEvents: []learnValue{
			{
				req:      Learn{From: 2, Rnd: 1, Val: "Lamport"},
				wantResp: ZeroValue,
			},
			{
				req:      Learn{From: 1, Rnd: 1, Val: "Leslie"},
				wantResp: ZeroValue,
			},
		},
	},
	{
		name:    "QuorumDifferentRounds",
		desc:    "3 nodes, single learn with rnd 2, then two learns with rnd 4 (quorum) -> report output and value of quorum",
		learner: NewLearner(0, 3),
		learnerEvents: []learnValue{
			{
				req:      Learn{From: 2, Rnd: 2, Val: "Lamport"},
				wantResp: ZeroValue,
			},
			{
				req:      Learn{From: 1, Rnd: 4, Val: "Leslie"},
				wantResp: ZeroValue,
			},
			{
				req:      Learn{From: 0, Rnd: 4, Val: "Leslie"},
				wantResp: "Leslie",
			},
		},
	},
}
