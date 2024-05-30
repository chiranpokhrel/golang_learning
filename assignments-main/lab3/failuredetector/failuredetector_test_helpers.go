// Test helper functions - DO NOT EDIT

package failuredetector

import (
	"slices"
	"sort"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/exp/maps"
)

const (
	ourID = 2
	delta = time.Second
)

var clusterOfThree = []int{2, 1, 0}

// setTestHook sets the test hook for the given failure detector and
// returns a channel that signals when the test hook has been called.
// The test hook is called at the start of the failure detector's main loop
// to reset the timeout signal's channel to nil. This is avoid having to
// wait for actual timeouts when testing the failure detector's timeout.
//
// Tests using the test hook must do <-done to wait for the test hook to be
// called between each fd.Deliver(heartbeat) and <-fd.Heartbeats() call.
func setTestHook(fd *EvtFailureDetector) <-chan bool {
	done := make(chan bool)
	fd.testHook = func() {
		fd.timeoutSignal.C = nil
		done <- true
	}
	return done
}

var transformer = cmpopts.AcyclicTransformer("SortedSlice", func(in interface{}) []int {
	switch v := in.(type) {
	case []int:
		out := slices.Clone(v) // Clone the slice to avoid mutating the original
		sort.Ints(out)         // Sort the slice to ensure consistent order for comparison
		return out
	case map[int]bool:
		keys := maps.Keys(v)
		sort.Ints(keys) // Sort the keys to ensure consistent order for comparison
		return keys
	default:
		panic("unsupported type")
	}
})

func createHeartbeats(destinations []int) []Heartbeat {
	heartbeats := make([]Heartbeat, len(destinations))
	for i, to := range destinations {
		heartbeats[i] = Heartbeat{To: to, From: ourID, Type: hbRequest}
	}
	return heartbeats
}

// Accumulator is simply a implementation of the SuspectRestorer interface
// that record the suspect and restore indications it receives. Used for
// testing.
type Accumulator struct {
	Suspects map[int]bool
	Restores map[int]bool
}

func NewAccumulator() *Accumulator {
	return &Accumulator{
		Suspects: make(map[int]bool),
		Restores: make(map[int]bool),
	}
}

func (a *Accumulator) Suspect(id int) {
	a.Suspects[id] = true
}

func (a *Accumulator) Restore(id int) {
	a.Restores[id] = true
}

func (a *Accumulator) Reset() {
	a.Suspects = make(map[int]bool)
	a.Restores = make(map[int]bool)
}
