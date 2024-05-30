package multipaxos

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMultiPaxosAcceptor(t *testing.T) {
	ignorePValueOrder := cmpopts.SortSlices(func(a, b PValue) bool {
		return a.Slot < b.Slot
	})
	for _, test := range acceptorTests {
		t.Run(test.name, func(t *testing.T) {
			acceptor := NewAcceptor(0)
			for _, msg := range test.acceptorMsgs {
				switch msg := msg.(type) {
				case preparePromise:
					prepare, wantPromise := msg.prepare, msg.wantPromise
					gotPromise := acceptor.handlePrepare(prepare)
					if diff := cmp.Diff(wantPromise, gotPromise, ignorePValueOrder); diff != "" {
						t.Log(msg.desc)
						t.Errorf("handlePrepare() mismatch (-want +got):\n%s", diff)
					}
				case acceptLearn:
					accept, wantLearn := msg.accept, msg.wantLearn
					gotLearn := acceptor.handleAccept(accept)
					if diff := cmp.Diff(wantLearn, gotLearn); diff != "" {
						t.Log(msg.desc)
						t.Errorf("handleAccept() mismatch (-want +got):\n%s", diff)
					}
				default:
					t.Fatalf("unknown message type: %T", msg)
				}
			}
		})
	}
}

var acceptorTests = []struct {
	name         string
	acceptorMsgs []any
}{
	{
		name: "PrepareNoPriorHistory",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 0, initial acceptor round is NoRound (-1) -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 0},
				wantPromise: Promise{To: 2, From: 0, Rnd: 0},
			},
		},
	},
	{
		name: "PrepareSlot",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
		},
	},
	{
		name: "PrepareIgnoreLowerCrnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
			preparePromise{
				desc:        "prepare slot 1 with round 1 -> ignore due to lower crnd",
				prepare:     Prepare{From: 1, Slot: 1, Crnd: 1},
				wantPromise: Promise{},
			},
		},
	},
	{
		name: "AcceptSlotIgnorePrepare",
		acceptorMsgs: []any{
			acceptLearn{
				desc:      "accept slot 1 with round 2, current acceptor rnd should be NoRound (-1) -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 2, Val: valThree},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 2, Val: valThree},
			},
			preparePromise{
				desc:        "prepare slot 1 with round 1, previous seen accept/learn (rnd 2) -> ignore due to lower crnd",
				prepare:     Prepare{From: 1, Slot: 1, Crnd: 1},
				wantPromise: Promise{},
			},
		},
	},
	{
		name: "PrepareNoHistory",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 0 with round 2, no previous prepare or accepts -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 0, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
		},
	},
	{
		name: "TwoPreparesNoHistory",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 1 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 1, Slot: 1, Crnd: 1},
				wantPromise: Promise{To: 1, From: 0, Rnd: 1},
			},
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
		},
	},
	{
		name: "PrepareHigherCrnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with current round -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 2, Val: valOne},
			},
			preparePromise{
				desc:        "new prepare slot 1 with higher crnd -> promise with correct rnd and history (slot 1)",
				prepare:     Prepare{From: 1, Slot: 1, Crnd: 3},
				wantPromise: Promise{To: 1, From: 0, Rnd: 3, Accepted: []PValue{{Slot: 1, Vrnd: 2, Vval: valOne}}},
			},
		},
	},
	{
		name: "AcceptTwoSlotsHigherRnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with current round -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 1 with higher round -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 5, Val: valTwo},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 5, Val: valTwo},
			},
			acceptLearn{
				desc:      "accept slot 2 with higher round -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 2, Rnd: 8, Val: valThree},
				wantLearn: Learn{From: 0, Slot: 2, Rnd: 8, Val: valThree},
			},
		},
	},
	{
		name: "PreparePromiseAcceptLowerRnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 4 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 1, Slot: 1, Crnd: 4},
				wantPromise: Promise{To: 1, From: 0, Rnd: 4},
			},
			acceptLearn{
				desc:      "accept slot 1 with lower round (2) than the current round (4) -> ignore accept",
				accept:    Accept{From: 2, Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: Learn{},
			},
		},
	},
	{
		name: "NewPreparePromiseWithHistory",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 2 with round 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 3, Rnd: 2, Val: valTwo},
				wantLearn: Learn{From: 0, Slot: 3, Rnd: 2, Val: valTwo},
			},
			acceptLearn{
				desc:      "accept slot 4 with round 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 4, Rnd: 2, Val: valThree},
				wantLearn: Learn{From: 0, Slot: 4, Rnd: 2, Val: valThree},
			},
			preparePromise{
				desc:        "new prepare slot 2 with round 3 -> output promise with correct rnd and history (slot 3 and 4)",
				prepare:     Prepare{From: 1, Slot: 2, Crnd: 3},
				wantPromise: Promise{To: 1, From: 0, Rnd: 3, Accepted: []PValue{{Slot: 3, Vrnd: 2, Vval: valTwo}, {Slot: 4, Vrnd: 2, Vval: valThree}}},
			},
		},
	},
	{
		name: "NewPreparePromiseWithHistoryII",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 3 with different sender and lower round -> ignore accept",
				accept:    Accept{From: 1, Slot: 3, Rnd: 1, Val: valTwo},
				wantLearn: Learn{},
			},
			acceptLearn{
				desc:      "accept slot 4 with round 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 4, Rnd: 2, Val: valThree},
				wantLearn: Learn{From: 0, Slot: 4, Rnd: 2, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 5 with higher round (5) -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 5, Rnd: 5, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 5, Rnd: 5, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 7 with lower round (1) -> ignore accept",
				accept:    Accept{From: 1, Slot: 7, Rnd: 1, Val: valTwo},
				wantLearn: Learn{},
			},
			preparePromise{
				desc:        "new prepare slot 2 and higher round (7) -> output promise with correct rnd (7) and history (slot 4 and 5)",
				prepare:     Prepare{From: 1, Slot: 2, Crnd: 7},
				wantPromise: Promise{To: 1, From: 0, Rnd: 7, Accepted: []PValue{{Slot: 4, Vrnd: 2, Vval: valThree}, {Slot: 5, Vrnd: 5, Vval: valOne}}},
			},
		},
	},
	{
		name: "NewPreparePromiseWithHistoryIII",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 3 with different sender and lower round -> ignore accept",
				accept:    Accept{From: 1, Slot: 3, Rnd: 1, Val: valTwo},
				wantLearn: Learn{},
			},
			acceptLearn{
				desc:      "accept slot 4 with round 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 4, Rnd: 2, Val: valThree},
				wantLearn: Learn{From: 0, Slot: 4, Rnd: 2, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 6 with higher round (5) -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 6, Rnd: 5, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 6, Rnd: 5, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 7 with lower round (1) -> ignore accept",
				accept:    Accept{From: 1, Slot: 7, Rnd: 1, Val: valTwo},
				wantLearn: Learn{},
			},
			preparePromise{
				desc:        "new prepare slot 2 and higher round (7) -> output promise with correct rnd (7) and history (slot 4 and 6)",
				prepare:     Prepare{From: 1, Slot: 2, Crnd: 7},
				wantPromise: Promise{To: 1, From: 0, Rnd: 7, Accepted: []PValue{{Slot: 4, Vrnd: 2, Vval: valThree}, {Slot: 6, Vrnd: 5, Vval: valOne}}},
			},
		},
	},
	{
		name: "NewPreparePromiseWithHistoryIV",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 1 again but with higher round (5) -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 5, Val: valTwo},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 5, Val: valTwo},
			},
			acceptLearn{
				desc:      "accept slot 3 with different sender and lower round (1) -> ignore accept",
				accept:    Accept{From: 1, Slot: 3, Rnd: 1, Val: valThree},
				wantLearn: Learn{},
			},
			acceptLearn{
				desc:      "accept slot 4 with round 5 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 4, Rnd: 5, Val: valThree},
				wantLearn: Learn{From: 0, Slot: 4, Rnd: 5, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 4 again but with higher round (8) -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 4, Rnd: 8, Val: valThree},
				wantLearn: Learn{From: 0, Slot: 4, Rnd: 8, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 6 with higher round (11) -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 6, Rnd: 11, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 6, Rnd: 11, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 7 with lower round (1) -> ignore accept",
				accept:    Accept{From: 1, Slot: 7, Rnd: 1, Val: valTwo},
				wantLearn: Learn{},
			},
			preparePromise{
				desc:        "new prepare slot 2 and higher round (13) from 1 -> output promise with correct rnd (13) and history (slot 4 and 6)",
				prepare:     Prepare{From: 1, Slot: 2, Crnd: 13},
				wantPromise: Promise{To: 1, From: 0, Rnd: 13, Accepted: []PValue{{Slot: 4, Vrnd: 8, Vval: valThree}, {Slot: 6, Vrnd: 11, Vval: valOne}}},
			},
		},
	},
	{
		name: "AcceptLearnWithHigherRnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     Prepare{From: 2, Slot: 1, Crnd: 2},
				wantPromise: Promise{To: 2, From: 0, Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 2 with higher round (5) -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 2, Rnd: 5, Val: valThree},
				wantLearn: Learn{From: 0, Slot: 2, Rnd: 5, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 1 again with higher round after we previously sent accept slot 2 -> learn with correct slot, rnd, and value",
				accept:    Accept{From: 2, Slot: 1, Rnd: 8, Val: valTwo},
				wantLearn: Learn{From: 0, Slot: 1, Rnd: 8, Val: valTwo},
			},
		},
	},
}
