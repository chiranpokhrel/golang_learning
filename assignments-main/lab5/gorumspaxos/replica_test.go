package gorumspaxos

import "testing"

func TestClientRequestCommit(t *testing.T) {
	testClientRequestCommit(t, func() {})
}
