package multipaxos

// promiseAccepts is a helper struct for testing the proposer.
type promiseAccepts struct {
	desc        string
	promise     Promise
	wantAccepts []Accept
}

// preparePromise is a helper struct for testing the acceptor.
type preparePromise struct {
	desc        string
	prepare     Prepare
	wantPromise Promise
}

// acceptLearn is a helper struct for testing the acceptor.
type acceptLearn struct {
	desc      string
	accept    Accept
	wantLearn Learn
}

// learnValue is a helper struct for testing the learner.
type learnValue struct {
	learn     Learn
	wantValue Value
	wantSlot  Slot
}

// mockLD is a mock leader detector.
type mockLD struct{}

func (l *mockLD) Leader() int {
	return -1
}

func (l *mockLD) Subscribe() <-chan int {
	return make(chan int)
}

var (
	valOne = Value{
		ClientID:  "1234",
		ClientSeq: 42,
		Command:   "ls",
	}
	valTwo = Value{
		ClientID:  "5678",
		ClientSeq: 99,
		Command:   "rm",
	}
	valThree = Value{
		ClientID:  "1369",
		ClientSeq: 4,
		Command:   "mkdir",
	}
)
