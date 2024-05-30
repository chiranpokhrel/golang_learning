package gorumspaxos

import "testing"

func TestFiveReplicas(t *testing.T) {
	testFiveReplicas(t, func() {})
}

func TestLeaderFailure(t *testing.T) {
	testLeaderFailure(t, func() {})
}
