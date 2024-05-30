// Failure detector tests - DO NOT EDIT

package failuredetector

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAllNodesShouldBeAlivePreStart(t *testing.T) {
	acc := NewAccumulator()
	fd := NewEvtFailureDetector(ourID, clusterOfThree, acc, delta)
	done := setTestHook(fd)

	fd.Start()
	<-done

	if len(fd.alive) != len(clusterOfThree) {
		t.Errorf("alive set contains %d node ids, want %d", len(fd.alive), len(clusterOfThree))
	}

	for _, id := range clusterOfThree {
		alive := fd.alive[id]
		if !alive {
			t.Errorf("node %d was not set alive", id)
			continue
		}
	}
}

func TestSendReplyToHeartbeatRequest(t *testing.T) {
	acc := NewAccumulator()
	fd := NewEvtFailureDetector(ourID, clusterOfThree, acc, delta)
	done := setTestHook(fd)

	fd.Start()
	<-done

	for i := 0; i < 10; i++ {
		hbReq := Heartbeat{To: ourID, From: i, Type: hbRequest}
		fd.Deliver(hbReq)
		<-done
		select {
		case hbResp := <-fd.Heartbeats():
			if hbResp.To != i {
				t.Errorf("Want heartbeat response to id %d, got id %d", i, hbResp.To)
			}
			if hbResp.From != ourID {
				t.Errorf("Want heartbeat response from id %d, got id %d", ourID, hbResp.From)
			}
			if hbResp.Type {
				t.Errorf("Want heartbeat of type response, got %v", hbResp)
			}
		default:
			t.Errorf("expected heartbeat response from %d", i)
		}
	}
}

func TestSetAliveDueToHeartbeatReply(t *testing.T) {
	acc := NewAccumulator()
	fd := NewEvtFailureDetector(ourID, clusterOfThree, acc, delta)
	done := setTestHook(fd)

	fd.Start()
	<-done

	for i := 0; i < 15; i++ {
		hbReply := Heartbeat{To: ourID, From: i, Type: hbReply}
		fd.Deliver(hbReply)
		<-done
		select {
		case hb := <-fd.Heartbeats():
			t.Errorf("want no outgoing heartbeat, got %v", hb)
		default:
		}
		alive := fd.alive[hbReply.From]
		if !alive {
			t.Errorf("got heartbeat reply from %d, but node was not marked as alive", i)
		}
	}
}

var timeoutTests = []struct {
	name              string
	alive             map[int]bool
	suspected         map[int]bool
	wantPostSuspected map[int]bool
	wantSuspects      []int
	wantRestores      []int
	wantDelay         time.Duration
}{
	{
		name:              "All nodes alive",
		alive:             map[int]bool{2: true, 1: true, 0: true},
		suspected:         map[int]bool{},
		wantPostSuspected: map[int]bool{},
		wantSuspects:      []int{},
		wantRestores:      []int{},
		wantDelay:         delta,
	},
	{
		name:              "All nodes alive, node 0 suspected and restored",
		alive:             map[int]bool{2: true, 1: true, 0: true},
		suspected:         map[int]bool{0: true},
		wantPostSuspected: map[int]bool{},
		wantSuspects:      []int{},
		wantRestores:      []int{0},
		wantDelay:         delta + delta,
	},
	{
		name:              "All nodes alive, node 1 suspected and restored",
		alive:             map[int]bool{2: true, 1: true, 0: true},
		suspected:         map[int]bool{1: true},
		wantPostSuspected: map[int]bool{},
		wantSuspects:      []int{},
		wantRestores:      []int{1},
		wantDelay:         delta + delta,
	},
	{
		name:              "All nodes alive, node 2 suspected and restored",
		alive:             map[int]bool{2: true, 1: true, 0: true},
		suspected:         map[int]bool{2: true},
		wantPostSuspected: map[int]bool{},
		wantSuspects:      []int{},
		wantRestores:      []int{2},
		wantDelay:         delta + delta,
	},
	{
		name:              "All nodes alive, suspected and restored",
		alive:             map[int]bool{2: true, 1: true, 0: true},
		suspected:         map[int]bool{2: true, 1: true, 0: true},
		wantPostSuspected: map[int]bool{},
		wantSuspects:      []int{},
		wantRestores:      []int{2, 1, 0},
		wantDelay:         delta + delta,
	},
	{
		name:              "Nodes 0 and 1 alive, node 2 suspected",
		alive:             map[int]bool{1: true, 0: true},
		suspected:         map[int]bool{2: true},
		wantPostSuspected: map[int]bool{2: true},
		wantSuspects:      []int{},
		wantRestores:      []int{},
		wantDelay:         delta,
	},
	{
		name:              "Nodes 0 and 2 alive, node 1 suspected",
		alive:             map[int]bool{2: true, 0: true},
		suspected:         map[int]bool{1: true},
		wantPostSuspected: map[int]bool{1: true},
		wantSuspects:      []int{},
		wantRestores:      []int{},
		wantDelay:         delta,
	},
	{
		name:              "Nodes 1 and 2 alive, node 0 suspected",
		alive:             map[int]bool{2: true, 1: true},
		suspected:         map[int]bool{0: true},
		wantPostSuspected: map[int]bool{0: true},
		wantSuspects:      []int{},
		wantRestores:      []int{},
		wantDelay:         delta,
	},
	{
		name:              "Node 2 alive, nodes 0 and 1 suspected",
		alive:             map[int]bool{2: true},
		suspected:         map[int]bool{},
		wantPostSuspected: map[int]bool{1: true, 0: true},
		wantSuspects:      []int{1, 0},
		wantRestores:      []int{},
		wantDelay:         delta,
	},
	{
		name:              "Node 1 alive, nodes 0 and 2 suspected",
		alive:             map[int]bool{1: true},
		suspected:         map[int]bool{},
		wantPostSuspected: map[int]bool{2: true, 0: true},
		wantSuspects:      []int{2, 0},
		wantRestores:      []int{},
		wantDelay:         delta,
	},
	{
		name:              "Node 0 alive, nodes 1 and 2 suspected",
		alive:             map[int]bool{0: true},
		suspected:         map[int]bool{},
		wantPostSuspected: map[int]bool{2: true, 1: true},
		wantSuspects:      []int{2, 1},
		wantRestores:      []int{},
		wantDelay:         delta,
	},
	{
		name:              "All nodes suspected",
		alive:             map[int]bool{},
		suspected:         map[int]bool{},
		wantPostSuspected: map[int]bool{2: true, 1: true, 0: true},
		wantSuspects:      []int{2, 1, 0},
		wantRestores:      []int{},
		wantDelay:         delta,
	},
	{
		name:              "All nodes suspected II",
		alive:             map[int]bool{},
		suspected:         map[int]bool{2: true, 1: true, 0: true},
		wantPostSuspected: map[int]bool{2: true, 1: true, 0: true},
		wantSuspects:      []int{},
		wantRestores:      []int{},
		wantDelay:         delta,
	},
}

func TestTimeoutProcedure(t *testing.T) {
	for _, test := range timeoutTests {
		t.Run(test.name, func(t *testing.T) {
			acc := NewAccumulator()
			fd := NewEvtFailureDetector(ourID, clusterOfThree, acc, delta)
			done := setTestHook(fd)

			// Wait until blocked
			fd.Start()
			defer fd.Stop()
			<-done

			// Set our test data
			fd.alive = test.alive
			fd.suspected = test.suspected

			// Trigger timeout procedure
			fd.timeout()

			// Alive set should always be empty
			if len(fd.alive) > 0 {
				t.Errorf("len(fd.alive) = %d, want 0; the alive set should always be empty after timeout procedure completes", len(fd.alive))
			}

			if diff := cmp.Diff(test.wantPostSuspected, fd.suspected); diff != "" {
				t.Errorf("Unexpected 'suspected' set after timeout procedure: (-want +got):\n%s", diff)
			}

			if fd.delay != test.wantDelay {
				t.Errorf("fd.delay = %v after timeout procedure, want %v", fd.delay, test.wantDelay)
			}

			// Check for the suspects we want
			if diff := cmp.Diff(test.wantSuspects, acc.Suspects, transformer); diff != "" {
				t.Errorf("Set of suspect indications differ: (-want +got):\n%s", diff)
			}
			// Check for the restores we want
			if diff := cmp.Diff(test.wantRestores, acc.Restores, transformer); diff != "" {
				t.Errorf("Set of restore indications differ: (-want +got):\n%s", diff)
			}

			// Check outgoing heartbeat requests
			var (
				gotHeartbeats  []Heartbeat
				wantHeartbeats = createHeartbeats(clusterOfThree)
			)

		hbReqCollect:
			for {
				select {
				case hbReq := <-fd.Heartbeats():
					gotHeartbeats = append(gotHeartbeats, hbReq)
				default:
					break hbReqCollect
				}
			}

			if diff := cmp.Diff(wantHeartbeats, gotHeartbeats); diff != "" {
				t.Errorf("Set of outgoing heartbeat requests differ: (-want +got):\n%s", diff)
			}
		})
	}
}
