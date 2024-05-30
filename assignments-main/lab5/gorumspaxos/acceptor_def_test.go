package gorumspaxos

import (
	"testing"

	pb "dat520/lab5/gorumspaxos/proto"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func testAcceptor(t *testing.T, handler func()) {
	sortAcceptedField := protocmp.SortRepeatedFields(&pb.PromiseMsg{}, "Accepted")
	for _, test := range acceptorTests {
		if ok := t.Run(test.name, func(t *testing.T) {
			acceptor := NewAcceptor()
			for _, msg := range test.acceptorMsgs {
				switch msg := msg.(type) {
				case preparePromise:
					prepare, wantPromise := msg.prepare, msg.wantPromise
					gotPromise := acceptor.handlePrepare(prepare)
					if diff := cmp.Diff(wantPromise, gotPromise, protocmp.Transform(), sortAcceptedField); diff != "" {
						t.Log(msg.desc)
						t.Errorf("handlePrepare() mismatch (-want +got):\n%s", diff)
					}
				case acceptLearn:
					accept, wantLearn := msg.accept, msg.wantLearn
					gotLearn := acceptor.handleAccept(accept)
					if diff := cmp.Diff(wantLearn, gotLearn, protocmp.Transform()); diff != "" {
						t.Log(msg.desc)
						t.Errorf("handleAccept() mismatch (-want +got):\n%s", diff)
					}
				default:
					t.Fatalf("unknown message type: %T", msg)
				}
			}
		}); ok {
			handler()
		}
	}
}

// preparePromise is a helper struct for testing the acceptor.
type preparePromise struct {
	desc        string
	prepare     *pb.PrepareMsg
	wantPromise *pb.PromiseMsg
}

// acceptLearn is a helper struct for testing the acceptor.
type acceptLearn struct {
	desc      string
	accept    *pb.AcceptMsg
	wantLearn *pb.LearnMsg
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
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 0},
				wantPromise: &pb.PromiseMsg{Rnd: 0},
			},
		},
	},
	{
		name: "PrepareSlot",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
		},
	},
	{
		name: "PrepareIgnoreLowerCrnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
			preparePromise{
				desc:        "prepare slot 1 with round 1 -> ignore due to lower crnd",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 1},
				wantPromise: nil,
			},
		},
	},
	{
		name: "AcceptSlotIgnorePrepare",
		acceptorMsgs: []any{
			acceptLearn{
				desc:      "accept slot 1 with round 2, current acceptor rnd should be NoRound (-1) -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 2, Val: valThree},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 2, Val: valThree},
			},
			preparePromise{
				desc:        "prepare slot 1 with round 1, previous seen accept/learn (rnd 2) -> ignore due to lower crnd",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 1},
				wantPromise: nil,
			},
		},
	},
	{
		name: "PrepareNoHistory",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 0 with round 2, no previous prepare or accepts -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 0, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
		},
	},
	{
		name: "TwoPreparesNoHistory",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 1 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 1},
				wantPromise: &pb.PromiseMsg{Rnd: 1},
			},
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
		},
	},
	{
		name: "PrepareHigherCrnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with current round -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 2, Val: valOne},
			},
			preparePromise{
				desc:        "new prepare slot 1 with higher crnd -> promise with correct rnd and history (slot 1)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 3},
				wantPromise: &pb.PromiseMsg{Rnd: 3, Accepted: []*pb.PValue{{Slot: 1, Vrnd: 2, Vval: valOne}}},
			},
		},
	},
	{
		name: "AcceptTwoSlotsHigherRnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with current round -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 1 with higher round -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 5, Val: valTwo},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 5, Val: valTwo},
			},
			acceptLearn{
				desc:      "accept slot 2 with higher round -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 2, Rnd: 8, Val: valThree},
				wantLearn: &pb.LearnMsg{Slot: 2, Rnd: 8, Val: valThree},
			},
		},
	},
	{
		name: "PreparePromiseAcceptLowerRnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 4 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 4},
				wantPromise: &pb.PromiseMsg{Rnd: 4},
			},
			acceptLearn{
				desc:      "accept slot 1 with lower round (2) than the current round (4) -> ignore accept",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: nil,
			},
		},
	},
	{
		name: "NewPreparePromiseWithHistory",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 2 with round 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 3, Rnd: 2, Val: valTwo},
				wantLearn: &pb.LearnMsg{Slot: 3, Rnd: 2, Val: valTwo},
			},
			acceptLearn{
				desc:      "accept slot 4 with round 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 4, Rnd: 2, Val: valThree},
				wantLearn: &pb.LearnMsg{Slot: 4, Rnd: 2, Val: valThree},
			},
			preparePromise{
				desc:        "new prepare slot 2 with round 3 -> output promise with correct rnd and history (slot 3 and 4)",
				prepare:     &pb.PrepareMsg{Slot: 2, Crnd: 3},
				wantPromise: &pb.PromiseMsg{Rnd: 3, Accepted: []*pb.PValue{{Slot: 3, Vrnd: 2, Vval: valTwo}, {Slot: 4, Vrnd: 2, Vval: valThree}}},
			},
		},
	},
	{
		name: "NewPreparePromiseWithHistoryII",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 3 with different sender and lower round -> ignore accept",
				accept:    &pb.AcceptMsg{Slot: 3, Rnd: 1, Val: valTwo},
				wantLearn: nil,
			},
			acceptLearn{
				desc:      "accept slot 4 with round 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 4, Rnd: 2, Val: valThree},
				wantLearn: &pb.LearnMsg{Slot: 4, Rnd: 2, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 5 with higher round (5) -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 5, Rnd: 5, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 5, Rnd: 5, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 7 with lower round (1) -> ignore accept",
				accept:    &pb.AcceptMsg{Slot: 7, Rnd: 1, Val: valTwo},
				wantLearn: nil,
			},
			preparePromise{
				desc:        "new prepare slot 2 and higher round (7) -> output promise with correct rnd (7) and history (slot 4 and 5)",
				prepare:     &pb.PrepareMsg{Slot: 2, Crnd: 7},
				wantPromise: &pb.PromiseMsg{Rnd: 7, Accepted: []*pb.PValue{{Slot: 4, Vrnd: 2, Vval: valThree}, {Slot: 5, Vrnd: 5, Vval: valOne}}},
			},
		},
	},
	{
		name: "NewPreparePromiseWithHistoryIII",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 3 with different sender and lower round -> ignore accept",
				accept:    &pb.AcceptMsg{Slot: 3, Rnd: 1, Val: valTwo},
				wantLearn: nil,
			},
			acceptLearn{
				desc:      "accept slot 4 with round 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 4, Rnd: 2, Val: valThree},
				wantLearn: &pb.LearnMsg{Slot: 4, Rnd: 2, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 6 with higher round (5) -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 6, Rnd: 5, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 6, Rnd: 5, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 7 with lower round (1) -> ignore accept",
				accept:    &pb.AcceptMsg{Slot: 7, Rnd: 1, Val: valTwo},
				wantLearn: nil,
			},
			preparePromise{
				desc:        "new prepare slot 2 and higher round (7) -> output promise with correct rnd (7) and history (slot 4 and 6)",
				prepare:     &pb.PrepareMsg{Slot: 2, Crnd: 7},
				wantPromise: &pb.PromiseMsg{Rnd: 7, Accepted: []*pb.PValue{{Slot: 4, Vrnd: 2, Vval: valThree}, {Slot: 6, Vrnd: 5, Vval: valOne}}},
			},
		},
	},
	{
		name: "NewPreparePromiseWithHistoryIV",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 1 again but with higher round (5) -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 5, Val: valTwo},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 5, Val: valTwo},
			},
			acceptLearn{
				desc:      "accept slot 3 with different sender and lower round (1) -> ignore accept",
				accept:    &pb.AcceptMsg{Slot: 3, Rnd: 1, Val: valThree},
				wantLearn: nil,
			},
			acceptLearn{
				desc:      "accept slot 4 with round 5 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 4, Rnd: 5, Val: valThree},
				wantLearn: &pb.LearnMsg{Slot: 4, Rnd: 5, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 4 again but with higher round (8) -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 4, Rnd: 8, Val: valThree},
				wantLearn: &pb.LearnMsg{Slot: 4, Rnd: 8, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 6 with higher round (11) -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 6, Rnd: 11, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 6, Rnd: 11, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 7 with lower round (1) -> ignore accept",
				accept:    &pb.AcceptMsg{Slot: 7, Rnd: 1, Val: valTwo},
				wantLearn: nil,
			},
			preparePromise{
				desc:        "new prepare slot 2 and higher round (13) from 1 -> output promise with correct rnd (13) and history (slot 4 and 6)",
				prepare:     &pb.PrepareMsg{Slot: 2, Crnd: 13},
				wantPromise: &pb.PromiseMsg{Rnd: 13, Accepted: []*pb.PValue{{Slot: 4, Vrnd: 8, Vval: valThree}, {Slot: 6, Vrnd: 11, Vval: valOne}}},
			},
		},
	},
	{
		name: "AcceptLearnWithHigherRnd",
		acceptorMsgs: []any{
			preparePromise{
				desc:        "prepare slot 1 with round 2 -> correct promise without previous history (pvalues)",
				prepare:     &pb.PrepareMsg{Slot: 1, Crnd: 2},
				wantPromise: &pb.PromiseMsg{Rnd: 2},
			},
			acceptLearn{
				desc:      "accept slot 1 with round 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 2, Val: valOne},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 2, Val: valOne},
			},
			acceptLearn{
				desc:      "accept slot 2 with higher round (5) -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 2, Rnd: 5, Val: valThree},
				wantLearn: &pb.LearnMsg{Slot: 2, Rnd: 5, Val: valThree},
			},
			acceptLearn{
				desc:      "accept slot 1 again with higher round after we previously sent accept slot 2 -> learn with correct slot, rnd, and value",
				accept:    &pb.AcceptMsg{Slot: 1, Rnd: 8, Val: valTwo},
				wantLearn: &pb.LearnMsg{Slot: 1, Rnd: 8, Val: valTwo},
			},
		},
	},
}
