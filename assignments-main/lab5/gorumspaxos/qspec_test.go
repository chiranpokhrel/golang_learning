package gorumspaxos

import (
	"testing"
)

func TestPaxosQSpec(t *testing.T) {
	testPaxosQSpec(t, func() {})
}

func TestPrepareQF(t *testing.T) {
	testPrepareQF(t, func() {})
}

func TestAcceptQF(t *testing.T) {
	testAcceptQF(t, func() {})
}

func TestClientHandleQF(t *testing.T) {
	testClientHandleQF(t, func() {})
}
