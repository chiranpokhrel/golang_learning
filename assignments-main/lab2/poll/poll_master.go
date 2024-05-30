package poll

// master for the poll.
type master struct {
	// TODO(student): Add necessary fields. Configuration, etc
}

// NewPollMaster creates a new poll master and starts the Gorums server.
func NewPollMaster(srvAddr string) (*master, error) {
	// TODO(student): Complete the method.
	// TODO(student): Remove comments below that are unnecessary when your implementation is complete.
	// TODO(student): Rephrase comments that help readability.
	// Complete the following steps:
	// 1. Create a gorums manager and use the manager to create an empty configuration.
	// 2. Store the gorums manager and configuration object in master struct.
	return &master{}, nil
}

// SrvAddress return the poll master's address for incoming requests.
func (pm *master) SrvAddress() string {
	return ""
}

// NumParticipants returns the number of participants currently in the configuration.
func (pm *master) NumParticipants() int {
	return 0
}

// HandleQuestion sends the question to all the participants.
func (pm *master) HandleQuestion(q *Question) (*Response, error) {
	// TODO(student): Complete the method.
	// TODO(student): Remove comments below that are unnecessary when your implementation is complete.
	// TODO(student): Rephrase comments that help readability.
	// Complete the following steps:
	// 1. Send the question to all the participants in the Gorums configuration.
	// 2. Wait for a response from a quorum of participants.
	return nil, nil
}

// TODO(student): Implement the PollMaster Service RPCs.

// TODO(student): Implement a QuorumSpec object.

// TODO(student): Implement the quorum function SendQuestionQF as a method on the QuorumSpec object to handle the responses from the participants.
// Complete the following steps:
// 1. SendQuestionQF returns the majority answer, if such an answer exists.
// 2. The majority answer is the valid answer that was answered by at least half of the participants.
// 3. If no majority answer exists, respond with a InconclusiveAnswer Response.
// 4. If not all participants respond within "QuestionTimeout" seconds, send an error response.
