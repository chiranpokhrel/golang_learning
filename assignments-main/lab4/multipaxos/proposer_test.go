package multipaxos

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMultiPaxosProposer(t *testing.T) {
	for _, test := range proposerTests {
		t.Run(test.name, func(t *testing.T) {
			for _, msg := range test.msgs {
				promise, wantAccepts := msg.promise, msg.wantAccepts
				gotAccepts := test.proposer.handlePromise(promise)
				if diff := cmp.Diff(wantAccepts, gotAccepts); diff != "" {
					t.Log(msg.desc)
					t.Errorf("handlePromise() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

var proposerTests = []struct {
	name     string
	proposer *Proposer
	msgs     []promiseAccepts
}{
	{
		name:     "SinglePromiseNoQuorum",
		proposer: NewProposer(2, 3, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "single promise from 1 with correct round, no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
		},
	},
	{
		name:     "TwoPromiseQuorum",
		proposer: NewProposer(2, 3, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with correct round, no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
			{
				desc:    "promise from 0, with correct round, quorum -> output and empty accept slice ",
				promise: Promise{To: 2, From: 0, Rnd: 2}, wantAccepts: []Accept{},
			},
		},
	},
	{
		name:     "PromiseDifferentRound",
		proposer: NewProposer(2, 3, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with different round (42) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 42}, wantAccepts: nil,
			},
		},
	},
	{
		name:     "PromiseDifferentRoundII",
		proposer: NewProposer(2, 3, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with different round (1) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 1}, wantAccepts: nil,
			},
			{
				desc:    "promise from 1 with different round (6) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 6}, wantAccepts: nil,
			},
			{
				desc:    "promise from 0 with different round (4) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 0, Rnd: 4}, wantAccepts: nil,
			},
		},
	},
	{
		name:     "PromiseSameRoundNoQuorum",
		proposer: NewProposer(2, 3, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with correct round (2), no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
			{
				desc:    "again promise from 1 with correct round (2), ignore, no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
			{
				desc:    "again promise from 1 with correct round (2), ignore, no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
		},
	},
	{
		name:     "PromiseSameRoundNoQuorumII",
		proposer: NewProposer(2, 3, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with different round (6) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 6}, wantAccepts: nil,
			},
			{
				desc:    "promise from 0 with different round (6) than proposer's (2), quorum for round 6, not but not our round, ignore -> no output",
				promise: Promise{To: 2, From: 0, Rnd: 6}, wantAccepts: nil,
			},
		},
	},
	{
		name:     "PValueScenario1",
		proposer: NewProposer(2, 3, 1, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "scenario 1 - message 1 - see figure in README.md",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
			{
				desc: "scenario 1 - message 2 - see figure in README.md",
				promise: Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 1, Vval: valOne},
				}},
				wantAccepts: []Accept{{From: 2, Slot: 2, Rnd: 2, Val: valOne}},
			},
		},
	},
	{
		name:     "PValueScenario2",
		proposer: NewProposer(2, 3, 1, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc: "scenario 2 - message 1 - see figure in README.md",
				promise: Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 0, Vval: valOne},
				}},
				wantAccepts: nil,
			},
			{
				desc: "scenario 2 - message 2 - see figure in README.md",
				promise: Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 1, Vval: valTwo},
				}},
				wantAccepts: []Accept{{From: 2, Slot: 2, Rnd: 2, Val: valTwo}},
			},
		},
	},
	{
		name:     "PValueScenario3",
		proposer: NewProposer(2, 3, 1, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc: "scenario 3 - message 1 - see figure in README.md",
				promise: Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
					{Slot: 1, Vrnd: 0, Vval: valTwo},
					{Slot: 2, Vrnd: 0, Vval: valOne},
				}},
				wantAccepts: nil,
			},
			{
				desc: "scenario 3 - message 2 - see figure in README.md",
				promise: Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 1, Vval: valTwo},
				}},
				wantAccepts: []Accept{{From: 2, Slot: 2, Rnd: 2, Val: valTwo}},
			},
		},
	},
	{
		name:     "PValueScenario4",
		proposer: NewProposer(2, 3, 1, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc: "scenario 4 - message 1 - see figure in README.md",
				promise: Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 1, Vval: valOne},
					{Slot: 4, Vrnd: 1, Vval: valThree},
				}},
				wantAccepts: nil,
			},
			{
				desc: "scenario 4 - message 2 - see figure in README.md",
				promise: Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 1, Vval: valOne},
				}},
				wantAccepts: []Accept{
					{From: 2, Slot: 2, Rnd: 2, Val: valOne},
					{From: 2, Slot: 3, Rnd: 2, Val: Value{}},
					{From: 2, Slot: 4, Rnd: 2, Val: valThree},
				},
			},
		},
	},
	{
		name:     "PValueScenario5",
		proposer: NewProposer(2, 3, 1, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc: "scenario 5 - message 1 - see figure in README.md",
				promise: Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 0, Vval: valOne},
					{Slot: 4, Vrnd: 0, Vval: valThree},
					{Slot: 5, Vrnd: 1, Vval: valTwo},
				}},
				wantAccepts: nil,
			},
			{
				desc: "scenario 5 - message 2 - see figure in README.md",
				promise: Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 0, Vval: valOne},
					{Slot: 4, Vrnd: 0, Vval: valThree},
					{Slot: 5, Vrnd: 1, Vval: valTwo},
				}},
				wantAccepts: nil,
			},
			{
				desc: "scenario 5 - message 3 - see figure in README.md",
				promise: Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 0, Vval: valOne},
					{Slot: 4, Vrnd: 1, Vval: valTwo},
				}},
				wantAccepts: []Accept{
					{From: 2, Slot: 2, Rnd: 2, Val: valOne},
					{From: 2, Slot: 3, Rnd: 2, Val: Value{}},
					{From: 2, Slot: 4, Rnd: 2, Val: valTwo},
					{From: 2, Slot: 5, Rnd: 2, Val: valTwo},
				},
			},
		},
	},
	{
		name:     "PValueScenario6",
		proposer: NewProposer(2, 3, 1, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc: "variation of scenario 5 - study test code for details",
				promise: Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 0, Vval: valOne},
					{Slot: 4, Vrnd: 1, Vval: valThree},
					{Slot: 5, Vrnd: 1, Vval: valTwo},
				}},
				wantAccepts: nil,
			},
			{
				desc: "variation of scenario 5 - study test code for details",
				promise: Promise{To: 2, From: 1, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 0, Vval: valOne},
					{Slot: 4, Vrnd: 1, Vval: valThree},
					{Slot: 5, Vrnd: 1, Vval: valTwo},
				}},
				wantAccepts: nil,
			},
			{
				desc: "variation of scenario 5 - study test code for details",
				promise: Promise{To: 2, From: 0, Rnd: 2, Accepted: []PValue{
					{Slot: 2, Vrnd: 0, Vval: valOne},
					{Slot: 4, Vrnd: 0, Vval: valTwo},
				}},
				wantAccepts: []Accept{
					{From: 2, Slot: 2, Rnd: 2, Val: valOne},
					{From: 2, Slot: 3, Rnd: 2, Val: Value{}},
					{From: 2, Slot: 4, Rnd: 2, Val: valThree},
					{From: 2, Slot: 5, Rnd: 2, Val: valTwo},
				},
			},
		},
	},
	{
		name:     "FiveNodeSinglePromiseNoQuorum",
		proposer: NewProposer(2, 5, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "single promise from 1 with correct round, n=5, no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
		},
	},
	{
		name:     "FiveNodeThreePromiseQuorum",
		proposer: NewProposer(2, 5, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with correct round, n=5, no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
			{
				desc:    "promise from 3 with correct round, n=5, no quorum -> no output",
				promise: Promise{To: 2, From: 3, Rnd: 2}, wantAccepts: nil,
			},
			{
				desc:    "promise from 4, with correct round, n=5, quorum -> output and empty accept slice",
				promise: Promise{To: 2, From: 4, Rnd: 2}, wantAccepts: []Accept{},
			},
		},
	},
	{
		name:     "FiveNodePromiseDifferentRound",
		proposer: NewProposer(2, 5, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with different round (42) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 42}, wantAccepts: nil,
			},
		},
	},
	{
		name:     "FiveNodePromiseDifferentRoundII",
		proposer: NewProposer(2, 5, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with different round (1) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 1}, wantAccepts: nil,
			},
			{
				desc:    "promise from 2 with different round (6) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 2, Rnd: 6}, wantAccepts: nil,
			},
			{
				desc:    "promise from 4 with different round (4) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 4, Rnd: 4}, wantAccepts: nil,
			},
		},
	},
	{
		name:     "FiveNodePromiseSameRoundNoQuorum",
		proposer: NewProposer(2, 5, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with correct round (2), no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
			{
				desc:    "again promise from 1 with correct round (2), ignore, no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
			{
				desc:    "again promise from 1 with correct round (2), ignore, no quorum -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 2}, wantAccepts: nil,
			},
		},
	},
	{
		name:     "FiveNodePromiseSameRoundNoQuorumII",
		proposer: NewProposer(2, 5, 0, &mockLD{}),
		msgs: []promiseAccepts{
			{
				desc:    "promise from 1 with different round (6) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 1, Rnd: 6}, wantAccepts: nil,
			},
			{
				desc:    "promise from 3 with different round (6) than proposer's (2), ignore -> no output",
				promise: Promise{To: 2, From: 3, Rnd: 6}, wantAccepts: nil,
			},
			{
				desc:    "promise from 4 with different round (6) than proposer's (2), quorum for round 6, not but not our round, ignore -> no output",
				promise: Promise{To: 2, From: 4, Rnd: 6}, wantAccepts: nil,
			},
		},
	},
}
