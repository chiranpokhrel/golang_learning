package singlepaxos

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSinglePaxosProposer(t *testing.T) {
	for _, test := range proposerTests {
		t.Run(test.name, func(t *testing.T) {
			test.proposer.clientValue = test.wantClientValue
			var gotAccept Accept
			for _, event := range test.proposerEvents {
				promise, wantAccept := event.req, event.wantResp
				gotAccept = test.proposer.handlePromise(promise)
				if diff := cmp.Diff(wantAccept, gotAccept); diff != "" {
					t.Log(test.desc)
					t.Errorf("handlePromise() mismatch (-want +got):\n%s", diff)
				}
				if test.proposer.clientValue != test.wantClientValue {
					t.Errorf("Unexpected Proposer.clientValue = %q, want %q", test.proposer.clientValue, test.wantClientValue)
				}
			}
			// Check if the final accept value is correct
			if gotAccept.Val != test.wantClientValue {
				t.Errorf("Accept.Val = %q, want %q", gotAccept.Val, test.wantClientValue)
			}
		})
	}
}

var proposerTests = []struct {
	name            string
	desc            string
	proposer        *Proposer
	proposerEvents  []promiseAccept
	wantClientValue Value
}{
	{
		name:     "NoQuorum",
		desc:     "no quorum -> no output",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 0, From: 1, Rnd: 0, Vrnd: NoRound, Vval: ZeroValue},
				wantResp: Accept{},
			},
		},
		wantClientValue: ZeroValue,
	},
	{
		name:     "FreeValueI",
		desc:     "valid quorum and no value reported -> propose (send accept) client value (free value) from proposer.clientValue field (I)",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 0, From: 1, Rnd: 2, Vrnd: NoRound, Vval: ZeroValue},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 0, Rnd: 2, Vrnd: NoRound, Vval: ZeroValue},
				wantResp: Accept{From: 2, Rnd: 2, Val: valueFromClientOne},
			},
		},
		wantClientValue: valueFromClientOne,
	},
	{
		name:     "FreeValueII",
		desc:     "valid quorum and no value reported -> propose (send accept) client value (free value) from proposer.clientValue field (II)",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 2, From: 1, Rnd: 2, Vrnd: NoRound, Vval: ZeroValue},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 0, Rnd: 2, Vrnd: NoRound, Vval: ZeroValue},
				wantResp: Accept{From: 2, Rnd: 2, Val: valueFromClientTwo},
			},
		},
		wantClientValue: valueFromClientTwo,
	},
	{
		name:     "CorrectValue",
		desc:     "valid quorum and a value reported -> propose correct value in accept",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 2, From: 1, Rnd: 2, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 0, Rnd: 2, Vrnd: 1, Vval: "Leslie"},
				wantResp: Accept{From: 2, Rnd: 2, Val: "Leslie"},
			},
		},
		wantClientValue: "Leslie",
	},
	{
		name:     "IgnoreDifferentRound",
		desc:     "promise for different round (1) than our current one (2) -> ignore promise",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 2, From: 1, Rnd: 1, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
		},
		wantClientValue: ZeroValue,
	},
	{
		name:     "IgnoreDifferentRoundII",
		desc:     "three promises, all for different rounds (1,6,4) than our current one (2) -> ignore all promises",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 2, From: 1, Rnd: 1, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 0, Rnd: 6, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 1, Rnd: 4, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
		},
		wantClientValue: ZeroValue,
	},
	{
		name:     "IgnoreIdenticalPromises",
		desc:     "three identical promises from the same sender -> no quorum, no output",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 2, From: 1, Rnd: 2, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 1, Rnd: 2, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 1, Rnd: 2, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
		},
		wantClientValue: ZeroValue,
	},
	{
		name: "IgnoreIdenticalPromisesII",
		desc: "three identical promises from the same sender for our round  -> no quorum, no output\n" +
			"then single promise for different round then ours -> ignore, no quorum, no output\n" +
			"then single promise for our round from last node -> quorum, report output and accept",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 2, From: 1, Rnd: 2, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 1, Rnd: 2, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 1, Rnd: 2, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 0, Rnd: 0, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 0, Rnd: 2, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{From: 2, Rnd: 2, Val: valueFromClientTwo},
			},
		},
		wantClientValue: valueFromClientTwo,
	},
	{
		name:     "DifferentValues",
		desc:     "valid quorum and two different values reported -> propose correct value (highest vrnd) in accept",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 2, From: 0, Rnd: 2, Vrnd: 1, Vval: "Lamport"},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 1, Rnd: 2, Vrnd: 0, Vval: "Leslie"},
				wantResp: Accept{From: 2, Rnd: 2, Val: "Lamport"},
			},
		},
		wantClientValue: "Lamport",
	},
	{
		name:     "DifferentRound",
		desc:     "two promises (majority) for different round than our current one -> ignore all promises and no output",
		proposer: NewProposer(2, 3),
		proposerEvents: []promiseAccept{
			{
				req:      Promise{To: 2, From: 0, Rnd: 6, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
			{
				req:      Promise{To: 2, From: 1, Rnd: 6, Vrnd: NoRound, Vval: ""},
				wantResp: Accept{},
			},
		},
		wantClientValue: ZeroValue,
	},
}

var incCrndTests = []struct {
	id    int
	n     int
	crnds []Round
}{
	{id: 0, n: 3, crnds: []Round{0, 3, 6, 9, 12, 15, 18, 21}},
	{id: 1, n: 5, crnds: []Round{1, 6, 11, 16, 21, 26, 31, 36}},
	{id: 4, n: 7, crnds: []Round{4, 11, 18, 25, 32, 39, 46}},
}

func TestSinglePaxosIncreaseCrnd(t *testing.T) {
	for _, test := range incCrndTests {
		proposer := NewProposer(test.id, test.n)
		for _, wantCrnd := range test.crnds {
			if proposer.crnd != wantCrnd {
				t.Errorf("Proposer[%d].crnd = %d, want %d", test.id, proposer.crnd, wantCrnd)
			}
			proposer.increaseCrnd()
		}
	}
}
