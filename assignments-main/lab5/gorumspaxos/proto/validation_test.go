package proto_test

import (
	"testing"

	pb "dat520/lab5/gorumspaxos/proto"
)

func TestIsValidPValue(t *testing.T) {
	tests := []struct {
		prepare *pb.PrepareMsg
		pval    *pb.PValue
		want    bool
	}{
		{prepare: nil, want: false},
		{pval: nil, want: false},
		{prepare: &pb.PrepareMsg{Slot: pb.NoSlot}, pval: &pb.PValue{Slot: 1, Vrnd: 1}, want: false},
		{prepare: &pb.PrepareMsg{Crnd: pb.NoRound}, pval: &pb.PValue{Slot: 1, Vrnd: 1}, want: false},
		{prepare: &pb.PrepareMsg{Slot: 1, Crnd: 2}, pval: &pb.PValue{Slot: pb.NoSlot}, want: false},
		{prepare: &pb.PrepareMsg{Slot: 1, Crnd: 2}, pval: &pb.PValue{Vrnd: pb.NoRound}, want: false},
		{prepare: &pb.PrepareMsg{Slot: 1, Crnd: 2}, pval: &pb.PValue{Slot: 1, Vrnd: 1}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 2, Crnd: 2}, pval: &pb.PValue{Slot: 1, Vrnd: 1}, want: false}, // pval's slot is old (< prepare.Slot)
		{prepare: &pb.PrepareMsg{Slot: 2, Crnd: 2}, pval: &pb.PValue{Slot: 2, Vrnd: 1}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 2, Crnd: 2}, pval: &pb.PValue{Slot: 3, Vrnd: 1}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 2, Crnd: 2}, pval: &pb.PValue{Slot: 9, Vrnd: 1}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 3, Crnd: 2}, pval: &pb.PValue{Slot: 3, Vrnd: 1}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 3, Crnd: 2}, pval: &pb.PValue{Slot: 4, Vrnd: 1}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 3, Crnd: 2}, pval: &pb.PValue{Slot: 9, Vrnd: 1}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 2}, pval: &pb.PValue{Slot: 5, Vrnd: 1}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 2}, pval: &pb.PValue{Slot: 5, Vrnd: 2}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 2}, pval: &pb.PValue{Slot: 5, Vrnd: 3}, want: false}, // acceptor has already accepted a value for a higher round than prepare.Crnd
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 2}, pval: &pb.PValue{Slot: 5, Vrnd: 4}, want: false}, // acceptor has already accepted a value for a higher round than prepare.Crnd
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 2}, pval: &pb.PValue{Slot: 5, Vrnd: 5}, want: false}, // acceptor has already accepted a value for a higher round than prepare.Crnd
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 9}, pval: &pb.PValue{Slot: 5, Vrnd: 9}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 9}, pval: &pb.PValue{Slot: 5, Vrnd: 8}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 9}, pval: &pb.PValue{Slot: 5, Vrnd: 7}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 9}, pval: &pb.PValue{Slot: 5, Vrnd: 6}, want: true},
		{prepare: &pb.PrepareMsg{Slot: 5, Crnd: 9}, pval: &pb.PValue{Slot: 5, Vrnd: 5}, want: true},
	}
	for _, tt := range tests {
		if got := tt.pval.IsValid(tt.prepare); got != tt.want {
			t.Errorf("pval(%v).IsValid(%v) = %v, want %v", tt.pval, tt.prepare, got, tt.want)
		}
	}
}
