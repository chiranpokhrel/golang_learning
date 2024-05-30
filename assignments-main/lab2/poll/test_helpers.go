package poll

import (
	"fmt"
	"time"
)

type testIn struct {
	Question             *Question
	participantResponses []*Response
}

func (ti testIn) String() string {
	out := fmt.Sprintf("Question: %+v\n\t\tParticipant Responses: [", ti.Question)
	for i, response := range ti.participantResponses {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%+v", response)
	}
	out += "]"
	return out
}

type testExpected struct {
	errIsNil bool
	response *Response
}

func (te testExpected) String() string {
	return fmt.Sprintf(
		`ErrIsNil: %v
		Response: %+v`,
		te.errIsNil, te.response)
}

type testCases struct {
	Description string
	In          testIn
	Expected    testExpected
}

func (tc testCases) String() string {
	return fmt.Sprintf(
		`
	Description: %v
	In:
		%v
	Expected:
		%v
	`, tc.Description, tc.In, tc.Expected,
	)
}

type ChannelQuestionHandler struct {
	answerChannel chan *Response
	questionOut   chan *Question
}

func (cqh *ChannelQuestionHandler) HandleQuestion(q *Question) (*Response, error) {
	cqh.questionOut <- q
	response := <-cqh.answerChannel
	return response, nil
}

type TimeoutQuestionHandler struct {
	wait time.Duration
}

func (tqh *TimeoutQuestionHandler) HandleQuestion(q *Question) (*Response, error) {
	// Wait until the request has timed out, before returning an answer
	<-time.After(tqh.wait)
	return &Response{
		ID:     q.ID,
		Answer: 0,
	}, nil
}

type AlwaysFirstQuestionHandler struct {
	wait time.Duration
}

func (AlwaysFirstQuestionHandler) HandleQuestion(q *Question) (*Response, error) {
	// Wait until the request has timed out, before returning an answer
	return &Response{
		ID:     q.ID,
		Answer: 0,
	}, nil
}
