package multipaxos

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMultiPaxosLearner(t *testing.T) {
	for _, test := range learnerTests {
		t.Run(test.name, func(t *testing.T) {
			for _, msg := range test.msgs {
				learn, wantValue, wantSlot := msg.learn, msg.wantValue, msg.wantSlot
				gotValue, gotSlot := test.learner.handleLearn(learn)
				if diff := cmp.Diff(wantValue, gotValue); diff != "" {
					t.Log(test.desc)
					t.Errorf("handleLearn() mismatch (-want +got):\n%s", diff)
				}
				if gotSlot != wantSlot {
					t.Log(test.desc)
					t.Errorf("handleLearn() slot id mismatch: got %v, want %v", gotSlot, msg.wantSlot)
				}
			}
		})
	}
}

var learnerTests = []struct {
	name    string
	desc    string
	learner *Learner
	msgs    []learnValue
}{
	{
		name:    "NoQuorum",
		desc:    "single learn for slot 1, 3 nodes, no quorum -> report zero value",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valOne}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "NoQuorumII",
		desc:    "single learn for slot 2, 3 nodes, no quorum -> report zero value",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 2, Rnd: 1, Val: valOne}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "TwoLearnsQuorum",
		desc:    "two learns for slot 1, 3 nodes, equal round and value, unique senders = quorum -> report output and value",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valOne}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 1, Rnd: 1, Val: valOne}, wantValue: valOne, wantSlot: 1},
		},
	},
	{
		name:    "TwoLearnsNoQuorum",
		desc:    "two learns for slot 1, 3 nodes, equal round and value, same sender = no quorum -> no output",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "TwoLearnsNoQuorumII",
		desc:    "two learns for slot 1, 3 nodes, same sender but different round = no quorum -> no output",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 2, Slot: 1, Rnd: 2, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 1, Rnd: 5, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "TwoLearnsNoQuorumIII",
		desc:    "two learns for slot 1, 3 nodes, different rounds, unique senders = no quorum -> no output",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 1, Rnd: 2, Val: valThree}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "TwoLearnsNoQuorumIV",
		desc:    "two learns for slot 1, 3 nodes, second learn should be ignored due to lower round -> no output",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 2, Slot: 1, Rnd: 2, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "NoQuorumAllSlots",
		desc:    "single learn for slot 1-6, 3 nodes, no quorum each slot -> no output for every learn received",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 2, Slot: 1, Rnd: 2, Val: valOne}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 2, Rnd: 1, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 3, Rnd: 5, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 4, Rnd: 2, Val: valOne}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 5, Rnd: 1, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 6, Rnd: 5, Val: valThree}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "QuorumAllSlots",
		desc:    "quorum of learns for slot 1-3, 3 nodes -> report output and value for every slot",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 0, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 1},
			{learn: Learn{From: 1, Slot: 2, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 0, Slot: 2, Rnd: 3, Val: valThree}, wantValue: valThree, wantSlot: 2},
			{learn: Learn{From: 0, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 3},
		},
	},
	{
		name:    "QuorumTwoLearns",
		desc:    "quorum of learns (two) for slot 1, the third learn -> report output and value but ignore last learn",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 0, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 1},
			{learn: Learn{From: 2, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "QuorumDifferentRounds",
		desc:    "single learn for slot 42, rnd 3, then quorum of learns (2) for slot 42 with higher rnd (4), -> report output and value after quorum",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 0, Slot: 42, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 42, Rnd: 4, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 42, Rnd: 4, Val: valThree}, wantValue: valThree, wantSlot: 42},
		},
	},
	{
		name:    "QuorumDifferentSlots",
		desc:    "quorum of learns for slot 1-3, 3 nodes, learns are not received in slot order, received in replica order, i.e. you need to store learns for different slots -> report output and value for every slot",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 0, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 0, Slot: 2, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 0, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 1},
			{learn: Learn{From: 1, Slot: 2, Rnd: 3, Val: valThree}, wantValue: valThree, wantSlot: 2},
			{learn: Learn{From: 1, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 3},
		},
	},
	{
		name:    "QuorumDifferentSlotsII",
		desc:    "quorum of learns for slot 3-1, 3 nodes, learns are received in mixed slot order (3,1,2) -> report output and value for every slot",
		learner: NewLearner(3),
		msgs: []learnValue{
			{learn: Learn{From: 0, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 3},
			{learn: Learn{From: 0, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 1},
			{learn: Learn{From: 1, Slot: 2, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 0, Slot: 2, Rnd: 3, Val: valThree}, wantValue: valThree, wantSlot: 2},
		},
	},
	{
		name:    "FiveNodeNoQuorum",
		desc:    "single learn for slot 1, 5 nodes, no quorum -> report zero value",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valOne}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "FiveNodeNoQuorumII",
		desc:    "single learn for slot 2, 5 nodes, no quorum -> report zero value",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 2, Rnd: 1, Val: valOne}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "FiveNodeThreeLearnsQuorum",
		desc:    "three learns for slot 1, 5 nodes, equal round and value, unique senders = quorum -> report output and value",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valOne}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 1, Rnd: 1, Val: valOne}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 3, Slot: 1, Rnd: 1, Val: valOne}, wantValue: valOne, wantSlot: 1},
		},
	},
	{
		name:    "FiveNodeThreeLearnsNoQuorum",
		desc:    "three learns for slot 1, 5 nodes, equal round and value, same sender = no quorum -> no output",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "FiveNodeThreeLearnsNoQuorumII",
		desc:    "three learns for slot 1, 5 nodes, same sender but different round = no quorum -> no output",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 2, Slot: 1, Rnd: 2, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 1, Rnd: 5, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 1, Rnd: 8, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "FiveNodeThreeLearnsNoQuorumIII",
		desc:    "three learns for slot 1, 5 nodes, different rounds, unique senders = no quorum -> no output",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 1, Rnd: 2, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 3, Slot: 1, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "FiveNodeThreeLearnsNoQuorumIV",
		desc:    "three learns for slot 1, 5 nodes, second and third learn should be ignored due to lower round -> no output",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 2, Slot: 1, Rnd: 2, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 1, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 3, Slot: 1, Rnd: 0, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "FiveNodeNoQuorumAllSlots",
		desc:    "single learn for slot 1-6, 5 nodes, no quorum each slot -> no output for every learn received",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 0, Slot: 1, Rnd: 1, Val: valOne}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 2, Rnd: 2, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 3, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 3, Slot: 4, Rnd: 4, Val: valOne}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 4, Slot: 5, Rnd: 5, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 6, Rnd: 6, Val: valThree}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "FiveNodeQuorumAllSlots",
		desc:    "quorum of learns for slot 1-3, 5 nodes -> report output and value for every slot",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 3, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 1},
			{learn: Learn{From: 1, Slot: 2, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 2, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 3, Slot: 2, Rnd: 3, Val: valThree}, wantValue: valThree, wantSlot: 2},
			{learn: Learn{From: 1, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 3, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 3},
		},
	},
	{
		name:    "FiveNodeQuorumFourLearns",
		desc:    "quorum of learns (three) for slot 1, the forth learn -> report output and value but ignore last learn",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 3, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 1},
			{learn: Learn{From: 4, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
		},
	},
	{
		name:    "FiveNodeQuorumDifferentRounds",
		desc:    "single learn for slot 42, rnd 3, then quorum of learns (2) for slot 42 with higher rnd (4), -> report output and value after quorum",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 4, Slot: 42, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 42, Rnd: 4, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 42, Rnd: 4, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 3, Slot: 42, Rnd: 4, Val: valThree}, wantValue: valThree, wantSlot: 42},
		},
	},
	{
		name:    "FiveNodeQuorumDifferentSlots",
		desc:    "quorum of learns for slot 1-3, 5 nodes, learns are not received in slot order, received in replica order, i.e. you need to store learns for different slots -> report output and value for every slot",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 0, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 0, Slot: 2, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 0, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 2, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 4, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 1},
			{learn: Learn{From: 4, Slot: 2, Rnd: 3, Val: valThree}, wantValue: valThree, wantSlot: 2},
			{learn: Learn{From: 4, Slot: 3, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 3},
		},
	},
	{
		name:    "FiveNodeQuorumDifferentSlotsII",
		desc:    "quorum of learns for slot 3-1, 5 nodes, learns are received in mixed slot order (3,1,2) -> report output and value for every slot",
		learner: NewLearner(5),
		msgs: []learnValue{
			{learn: Learn{From: 1, Slot: 3, Rnd: 3, Val: valOne}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 3, Rnd: 3, Val: valOne}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 4, Slot: 3, Rnd: 3, Val: valOne}, wantValue: valOne, wantSlot: 3},
			{learn: Learn{From: 3, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 0, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 1, Slot: 1, Rnd: 3, Val: valTwo}, wantValue: valTwo, wantSlot: 1},
			{learn: Learn{From: 1, Slot: 2, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 2, Slot: 2, Rnd: 3, Val: valThree}, wantValue: Value{}, wantSlot: 0},
			{learn: Learn{From: 0, Slot: 2, Rnd: 3, Val: valThree}, wantValue: valThree, wantSlot: 2},
		},
	},
}
