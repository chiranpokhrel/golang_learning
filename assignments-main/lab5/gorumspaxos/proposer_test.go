package gorumspaxos

import (
	"testing"
)

func TestPerformAccept(t *testing.T) {
	testPerformAccept(t, func() {})
}

func TestPerformCommit(t *testing.T) {
	testPerformCommit(t, func() {})
}

func TestRunPhaseOne(t *testing.T) {
	testRunPhaseOne(t, func() {})
}
