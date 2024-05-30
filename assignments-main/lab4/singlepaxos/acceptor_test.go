package singlepaxos

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSinglePaxosAcceptor(t *testing.T) {
	for _, test := range acceptorTests {
		t.Run(test.name, func(t *testing.T) {
			acceptor := NewAcceptor(0)
			for _, event := range test.acceptorEvents {
				switch event := event.(type) {
				case preparePromise:
					prepare, wantPromise := event.req, event.wantResp
					gotPromise := acceptor.handlePrepare(prepare)
					if diff := cmp.Diff(wantPromise, gotPromise); diff != "" {
						t.Log(test.desc)
						t.Errorf("handlePrepare() mismatch (-want +got):\n%s", diff)
					}
				case acceptLearn:
					accept, wantLearn := event.req, event.wantResp
					gotLearn := acceptor.handleAccept(accept)
					if diff := cmp.Diff(wantLearn, gotLearn); diff != "" {
						t.Log(test.desc)
						t.Errorf("handleAccept() mismatch (-want +got):\n%s", diff)
					}
				default:
					t.Log(test.desc)
					t.Fatalf("unknown event type: %T", event)
				}
			}
		})
	}
}

var acceptorTests = []struct {
	name           string
	desc           string
	acceptorEvents []any
}{
	{
		name: "SinglePrepare",
		desc: "no previous received prepare -> reply with correct rnd and no vrnd/vval",
		acceptorEvents: []any{
			preparePromise{
				req:      Prepare{From: 1, Crnd: 1},
				wantResp: Promise{To: 1, From: 0, Rnd: 1, Vrnd: NoRound, Vval: ZeroValue},
			},
		},
	},
	{
		name: "TwoPrepares",
		desc: "two prepares, the second with higher round -> reply correctly to both",
		acceptorEvents: []any{
			preparePromise{
				req:      Prepare{From: 1, Crnd: 1},
				wantResp: Promise{To: 1, From: 0, Rnd: 1, Vrnd: NoRound, Vval: ZeroValue},
			},
			preparePromise{
				req:      Prepare{From: 2, Crnd: 2},
				wantResp: Promise{To: 2, From: 0, Rnd: 2, Vrnd: NoRound, Vval: ZeroValue},
			},
		},
	},
	{
		name: "PrepareAccept",
		desc: "single prepare followed by corresponding accept -> emit learn. then new prepare with higher round -> report correct vrnd, vval",
		acceptorEvents: []any{
			preparePromise{
				req:      Prepare{From: 1, Crnd: 1},
				wantResp: Promise{To: 1, From: 0, Rnd: 1, Vrnd: NoRound, Vval: ZeroValue},
			},
			acceptLearn{
				req:      Accept{From: 1, Rnd: 1, Val: "Lamport"},
				wantResp: Learn{From: 0, Rnd: 1, Val: "Lamport"},
			},
			preparePromise{
				req:      Prepare{From: 2, Crnd: 2},
				wantResp: Promise{To: 2, From: 0, Rnd: 2, Vrnd: 1, Vval: "Lamport"},
			},
		},
	},
	{
		name: "PrepareLowerCrnd",
		desc: "prepare with crnd lower than seen rnd -> ignore prepare",
		acceptorEvents: []any{
			preparePromise{
				req:      Prepare{From: 1, Crnd: 2},
				wantResp: Promise{To: 1, From: 0, Rnd: 2, Vrnd: NoRound, Vval: ZeroValue},
			},
			preparePromise{
				req:      Prepare{From: 1, Crnd: 1},
				wantResp: Promise{},
			},
		},
	},
	{
		name: "AcceptLowerRnd",
		desc: "accept with lower rnd than what we have sent in promise -> ignore accept, i.e. no learn",
		acceptorEvents: []any{
			preparePromise{
				req:      Prepare{From: 1, Crnd: 2},
				wantResp: Promise{To: 1, From: 0, Rnd: 2, Vrnd: NoRound, Vval: ZeroValue},
			},
			acceptLearn{
				req:      Accept{From: 2, Rnd: 1, Val: "Lamport"},
				wantResp: Learn{},
			},
		},
	},
}
