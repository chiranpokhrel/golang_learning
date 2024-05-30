package poll

import (
	"errors"
	"time"
)

// This file contains constants and interface definitions used by the poll package.
// Please do not change the file

const (
	InconclusiveAnswer = -1

	GorumsDialTimeout = 1 * time.Second
	QuestionTimeout   = 10 * time.Second
	RegisterTimeout   = 5 * time.Second
)

var errNoParticipants = errors.New("no available participants")

type QuestionHandler interface {
	// HandleQuestion returns an answer to the provided question.
	HandleQuestion(*Question) (*Response, error)
}

type PollMaster interface {
	QuestionHandler
	SrvAddress() string
	NumParticipants() int
}

type PollParticipant interface {
	Register(string) error
	Deregister() error
}

type Question struct {
	ID       int      // ID of the question
	Question string   // The question itself
	Options  []string // The alternatives
}

type Response struct {
	ID     int // ID of the question to which this is a response
	Answer int // The index of the chosen option
}
