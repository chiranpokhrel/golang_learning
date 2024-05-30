package poll

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/relab/gorums"
)

// Tests that a participant can register and deregister from the poll master.
// This should be able to run without implementing the HandleQuestion() method.
func TestRegisterDeregister(t *testing.T) {
	pm, err := NewPollMaster("")
	if err != nil {
		t.Errorf("Unable to create PollMaster: %v", err)
	}

	p, err := NewParticipant(&AlwaysFirstQuestionHandler{}, "")
	if err != nil {
		t.Errorf("Unable to create participant: %v", err)
	}

	// Register the participant
	err = p.Register(pm.SrvAddress())
	if err != nil {
		t.Errorf("Unexpected error while registering participants: %v", err)
	}
	if pm.NumParticipants() != 1 {
		t.Errorf("Unexpected size of configuration. Expected 1 node, got %v nodes", pm.NumParticipants())
	}

	// Deregister the participant
	if err = p.Deregister(); err != nil {
		t.Errorf("Failed to deregister participants: %v", err)
	}
	if pm.NumParticipants() != 0 {
		t.Errorf("Unexpected size of configuration. Expected 0 node, got %v nodes", pm.NumParticipants())
	}
}

// Checks that NumParticipants returns the correct amount.
func TestNumParticipants(t *testing.T) {
	const (
		firstNumParticipants      = 10
		deregisterNumParticipants = 5
		SecondTestNumParticipants = firstNumParticipants - deregisterNumParticipants
		SecondAddNumParticipants  = 10
		ThirdTestNumParticipants  = SecondTestNumParticipants + SecondAddNumParticipants
	)

	pm, err := NewPollMaster("")
	if err != nil {
		t.Errorf("Unable to create PollMaster: %v", err)
	}

	participants := []PollParticipant{}
	for i := 0; i < firstNumParticipants; i++ {
		p, err := NewParticipant(&AlwaysFirstQuestionHandler{}, "")
		if err != nil {
			t.Errorf("Unable to create participant: %v", err)
		}
		p.Register(pm.SrvAddress())
		participants = append(participants, p)
	}

	if pm.NumParticipants() != firstNumParticipants {
		t.Errorf("Received wrong number of participants. Got %v. Expected 10", pm.NumParticipants())
	}

	for i := 0; i < deregisterNumParticipants; i++ {
		participants[i].Deregister()
	}

	if pm.NumParticipants() != SecondTestNumParticipants {
		t.Errorf("Received wrong number of participants. Got %v. Expected 5", pm.NumParticipants())
	}

	for i := 0; i < SecondAddNumParticipants; i++ {
		p, err := NewParticipant(&AlwaysFirstQuestionHandler{}, "")
		if err != nil {
			t.Errorf("Unable to create participant: %v", err)
		}
		p.Register(pm.SrvAddress())
		participants = append(participants, p)
	}

	if pm.NumParticipants() != ThirdTestNumParticipants {
		t.Errorf("Received wrong number of participants. Got %v. Expected 15", pm.NumParticipants())
	}

	for i := deregisterNumParticipants; i < len(participants); i++ {
		participants[i].Deregister()
	}

	if pm.NumParticipants() != 0 {
		t.Errorf("Received wrong number of participants. Got %v. Expected 0", pm.NumParticipants())
	}
}

// Tests that a participant starts receiving questions after registering
// Only works after HandleQuestions, Part 2, has been implemented
func TestRegisterReceiveQuestions(t *testing.T) {
	pm, err := NewPollMaster("")
	if err != nil {
		t.Errorf("Unable to create PollMaster: %v", err)
	}

	SentQuestion := Question{
		ID:       1,
		Question: "question 1",
		Options: []string{
			"Option 1",
			"Option 2",
		},
	}
	_, err = pm.HandleQuestion(&SentQuestion)
	if err != errNoParticipants {
		t.Errorf("No registered participants. Should receive 'noParticipantError'. Got: %v", err)
	}

	qhInChan := make(chan *Response, 1)
	qhQuestion := make(chan *Question, 1)
	qh := &ChannelQuestionHandler{qhInChan, qhQuestion}
	p, err := NewParticipant(qh, "")
	if err != nil {
		t.Errorf("Unable to create participants: %v", err)
	}
	err = p.Register(pm.SrvAddress())
	if err != nil {
		t.Errorf("Unexpected error while registering participants: %v", err)
	}
	defer p.Deregister()

	// Check that the participant receives a question after registering
	SentQuestion = Question{
		ID:       2,
		Question: "question 2",
		Options: []string{
			"Option 1",
			"Option 2",
		},
	}
	pm.HandleQuestion(&SentQuestion)

	select {
	// Check if a question is received within 5 seconds. If not fail the test
	case <-time.After(5 * time.Second):
		t.Errorf("The participant should begin receiving questions after registering")
	case q := <-qhQuestion:
		if !reflect.DeepEqual(q, &SentQuestion) {
			t.Errorf("The participant should receive the same question that was sent. Sent: %v, Got: %v", SentQuestion, q)
		}
	}
}

// Tests that a participant stops receiving questions after having deregistered.
// Only works after HandleQuestions, Part 2, has been implemented.
func TestDeregisterStopReceivingQuestions(t *testing.T) {
	pm, err := NewPollMaster("")
	if err != nil {
		t.Errorf("Unable to create PollMaster: %v", err)
	}

	// There needs to be at least one participant to send questions
	// So we create this to always respond
	p, err := NewParticipant(&AlwaysFirstQuestionHandler{}, "")
	if err != nil {
		t.Errorf("Unable to create participants: %v", err)
	}
	err = p.Register(pm.SrvAddress())
	if err != nil {
		t.Errorf("Unexpected error while registering participants: %v", err)
	}
	defer p.Deregister()

	qhInChan := make(chan *Response, 1)
	qhQuestion := make(chan *Question, 1)
	qh := &ChannelQuestionHandler{qhInChan, qhQuestion}
	p, err = NewParticipant(qh, "")
	if err != nil {
		t.Errorf("Unable to create participants: %v", err)
	}
	err = p.Register(pm.SrvAddress())
	if err != nil {
		t.Errorf("Unexpected error while registering participants: %v", err)
	}

	if err = p.Deregister(); err != nil {
		t.Errorf("Failed to deregister participants: %v", err)
	}
	SentQuestion := Question{
		ID:       1,
		Question: "question 1",
		Options: []string{
			"Option 1",
			"Option 2",
		},
	}
	pm.HandleQuestion(&SentQuestion)

	select {
	// Check if a question is received within 5 seconds. If it is fail the test
	case <-time.After(5 * time.Second):
	case <-qhQuestion:
		t.Errorf("The participant should not receive questions after having deregistered")
	}
}

// Tests that the service as a whole is working. mainly focuses on checking that the result is correct
func TestPollingService(t *testing.T) {
	pm, err := NewPollMaster("")
	if err != nil {
		t.Errorf("Unable to create PollMaster: %v", err)
	}

	// Start 4 participants
	qhChannels := []chan *Response{}
	qhQuestions := []chan *Question{}
	for i := 0; i < 4; i++ {
		// Create channels to send responses and receive questions on for each participant
		// These will be used to ensure that the questions are properly received
		// and to supply custom responses
		qhInChan := make(chan *Response, 1)
		qhQuestion := make(chan *Question, 1)

		qhChannels = append(qhChannels, qhInChan)
		qhQuestions = append(qhQuestions, qhQuestion)

		qh := &ChannelQuestionHandler{qhInChan, qhQuestion}
		p, err := NewParticipant(qh, "")
		if err != nil {
			t.Errorf("Unable to create participants: %v", err)
		}
		err = p.Register(pm.SrvAddress())
		if err != nil {
			t.Errorf("Unexpected error while registering participants: %v", err)
		}
		defer p.Deregister()
	}

	for _, testCase := range pollingServiceTests {
		// Send the respective answers to all question handlers
		for j, answer := range testCase.In.participantResponses {
			qhChannels[j] <- answer
		}
		resp, err := pm.HandleQuestion(testCase.In.Question)

		// If we don't expect an error check that we don't receive an error.
		// Otherwise check that we do receive an error
		if testCase.Expected.errIsNil && err != nil {
			t.Fatalf("Unexpected error: \n Got: %+v. \n%v", err, testCase)
		} else if !testCase.Expected.errIsNil && err == nil {
			t.Fatalf("Expected an error: \n Got: %v. \n%v", err, testCase)
		}

		// Check that the participant forwards the correct question
		for _, questionChan := range qhQuestions {
			receivedQuestion := <-questionChan
			if !reflect.DeepEqual(testCase.In.Question, receivedQuestion) {
				t.Fatalf("Did not receive correct question.\n Got: %+v. Expected: %+v%v", receivedQuestion, testCase.In.Question, testCase)
			}
		}
		// Check that the poll master arrives at the correct response.
		if !reflect.DeepEqual(resp, testCase.Expected.response) {
			t.Fatalf("Unexpected answer: \n Got: %+v. Expected: %+v \n%v", resp, testCase.Expected.response, testCase)
		}
	}
}

// Tests that the HandleQuestion times out after an expected time
func Test4Participants3Timeouts(t *testing.T) {
	pm, err := NewPollMaster("")
	if err != nil {
		t.Errorf("Unable to create PollMaster: %v", err)
	}

	for i := 0; i < 3; i++ {
		qh := &TimeoutQuestionHandler{QuestionTimeout}
		p, err := NewParticipant(qh, "")
		if err != nil {
			t.Errorf("Unable to create participants: %v", err)
		}
		err = p.Register(pm.SrvAddress())
		if err != nil {
			t.Errorf("Unexpected error while registering participants: %v", err)
		}
		defer p.Deregister()
	}
	qh := &AlwaysFirstQuestionHandler{}
	p, err := NewParticipant(qh, "")
	if err != nil {
		t.Errorf("Unable to create participants: %v", err)
	}
	p.Register(pm.SrvAddress())
	defer p.Deregister()

	resp, err := pm.HandleQuestion(&Question{
		ID:       0,
		Question: "Question 1",
		Options: []string{
			"Option 1",
			"Option 2",
			"Option 3",
			"Option 4",
		},
	})
	var expectedError gorums.QuorumCallError
	if !errors.As(err, &expectedError) {
		t.Errorf("Expected timeout error. Received: %v", err)
	}
	if resp != nil {
		t.Errorf("Unexpected response. On timeout we expect to receive nil. Got: %v", resp)
	}
}

var pollingServiceTests = []testCases{
	{
		Description: "All participants responds with same answer",
		In: testIn{
			Question: &Question{
				ID:       1,
				Question: "Question 1",
				Options: []string{
					"Option 1",
					"Option 2",
					"Option 3",
					"Option 4",
				},
			},
			participantResponses: []*Response{
				{1, 0}, {1, 0}, {1, 0}, {1, 0},
			},
		},
		Expected: testExpected{
			true,
			&Response{
				ID:     1,
				Answer: 0,
			},
		},
	},
	{
		Description: "All participants respond, no majority answer",
		In: testIn{
			Question: &Question{
				ID:       2,
				Question: "Question 2",
				Options: []string{
					"Option 1",
					"Option 2",
					"Option 3",
					"Option 4",
				},
			},
			participantResponses: []*Response{
				{2, 0}, {2, 0}, {2, 1}, {2, 3},
			},
		},
		Expected: testExpected{
			true,
			&Response{
				ID:     2,
				Answer: InconclusiveAnswer,
			},
		},
	},
	{
		Description: "3 values larger than valid options. Expect inconclusive answer",
		In: testIn{
			Question: &Question{
				ID:       3,
				Question: "Question 3",
				Options: []string{
					"Option 1",
					"Option 2",
					"Option 3",
					"Option 4",
				},
			},
			participantResponses: []*Response{
				{3, 0}, {3, 10}, {3, 10}, {3, 10},
			},
		},
		Expected: testExpected{
			true,
			&Response{
				ID:     3,
				Answer: InconclusiveAnswer,
			},
		},
	},
	{
		Description: "3 values smaller than valid options. Expect inconclusive answer",
		In: testIn{
			Question: &Question{
				ID:       4,
				Question: "Question 4",
				Options: []string{
					"Option 1",
					"Option 2",
					"Option 3",
					"Option 4",
				},
			},
			participantResponses: []*Response{
				{4, 0}, {4, -10}, {4, -10}, {4, -10},
			},
		},
		Expected: testExpected{
			true,
			&Response{
				ID:     4,
				Answer: InconclusiveAnswer,
			},
		},
	},
	{
		Description: "1 invalid value. Ignore it, select the other",
		In: testIn{
			Question: &Question{
				ID:       5,
				Question: "Question 5",
				Options: []string{
					"Option 1",
					"Option 2",
					"Option 3",
					"Option 4",
				},
			},
			participantResponses: []*Response{
				{5, 0}, {5, 0}, {5, 0}, {5, -10},
			},
		},
		Expected: testExpected{
			true,
			&Response{
				ID:     5,
				Answer: 0,
			},
		},
	},
	{
		Description: "3 values larger than valid options. Expect inconclusive answer",
		In: testIn{
			Question: &Question{
				ID:       6,
				Question: "Question 6",
				Options: []string{
					"Option 1",
					"Option 2",
					"Option 3",
					"Option 4",
				},
			},
			participantResponses: []*Response{
				{6, 4}, {6, 4}, {6, 4}, {6, 0},
			},
		},
		Expected: testExpected{
			true,
			&Response{
				ID:     6,
				Answer: InconclusiveAnswer,
			},
		},
	},
	{
		Description: "3 responses with wrong id. Expect inconclusive answer",
		In: testIn{
			Question: &Question{
				ID:       7,
				Question: "Question 7",
				Options: []string{
					"Option 1",
					"Option 2",
					"Option 3",
					"Option 4",
				},
			},
			participantResponses: []*Response{
				{5, 1}, {1, 1}, {1, 1}, {7, 1},
			},
		},
		Expected: testExpected{
			true,
			&Response{
				ID:     7,
				Answer: InconclusiveAnswer,
			},
		},
	},
	{
		Description: "3 of same answer, but 1 has wrong id. Expect inconclusive answer",
		In: testIn{
			Question: &Question{
				ID:       8,
				Question: "Question 8",
				Options: []string{
					"Option 1",
					"Option 2",
					"Option 3",
					"Option 4",
				},
			},
			participantResponses: []*Response{
				{8, 1}, {8, 1}, {1, 1}, {8, 2},
			},
		},
		Expected: testExpected{
			true,
			&Response{
				ID:     8,
				Answer: InconclusiveAnswer,
			},
		},
	},
}
