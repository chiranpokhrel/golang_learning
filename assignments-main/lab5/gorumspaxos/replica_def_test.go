package gorumspaxos

import (
	"testing"
	"time"

	pb "dat520/lab5/gorumspaxos/proto"

	"github.com/google/go-cmp/cmp"
	"github.com/relab/gorums"
	"google.golang.org/protobuf/testing/protocmp"
)

func testClientRequestCommit(t *testing.T, handler func()) {
	tests := []struct {
		desc     string
		slot     Slot
		val      *pb.Value
		wantResp *pb.Response
		wantErr  bool
	}{
		{desc: "Client request 1 commit", slot: 1, val: valOne, wantResp: &pb.Response{ClientID: valOne.ClientID, ClientSeq: valOne.ClientSeq, ClientCommand: valOne.ClientCommand}, wantErr: false},
		{desc: "Client request 2 commit", slot: 2, val: valTwo, wantResp: &pb.Response{ClientID: valTwo.ClientID, ClientSeq: valTwo.ClientSeq, ClientCommand: valTwo.ClientCommand}, wantErr: false},
		{desc: "Client request 3 commit", slot: 3, val: valThree, wantResp: &pb.Response{ClientID: valThree.ClientID, ClientSeq: valThree.ClientSeq, ClientCommand: valThree.ClientCommand}, wantErr: false},
		{desc: "Client request 4 commit", slot: 4, val: &pb.Value{}, wantResp: &pb.Response{}, wantErr: false},
	}
	replica := newTestReplicaLeader()
	for _, test := range tests {
		go func() {
			time.Sleep(1 * time.Millisecond)
			replica.Commit(gorums.ServerCtx{}, &pb.LearnMsg{Val: test.val, Slot: test.slot})
		}()
		resp, err := replica.ClientHandle(gorums.ServerCtx{}, test.val)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(resp, test.wantResp, protocmp.Transform()); diff != "" {
			t.Errorf("%s mismatch (-want +got):\n%s", t.Name(), diff)
		}
	}
	if !t.Failed() {
		handler()
	}
}
