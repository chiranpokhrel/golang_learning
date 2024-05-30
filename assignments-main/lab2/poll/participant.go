package poll

// participant in a poll.
type participant struct {
	qh QuestionHandler

	// TODO: Add more fields as needed
}

// NewParticipant creates the new poll participant and starts the Gorums server.
func NewParticipant(qh QuestionHandler, addr string) (*participant, error) {
	// TODO(student): Complete the method.
	// TODO(student): Remove comments below that are unnecessary when your implementation is complete.
	// TODO(student): Rephrase comments that help readability.
	// Complete the following steps:
	// 1. Start listening on the specified address.
	// 2. Create a gorums server and register participant as the service implementation.
	// 3. Start serving the gorums server using a goroutine.
	// 4. Store the local serving address and question handler in participant.
	return &participant{
		qh: qh,
	}, nil
}

// Register registers the participant to the given poll master.
func (p *participant) Register(srvAddr string) error {
	// TODO(student): Complete the method.
	// TODO(student): Remove comments below that are unnecessary when your implementation is complete.
	// TODO(student): Rephrase comments that help readability.
	// Complete the following steps:
	// 1. Create a connection to the master server and save the connection for future use.
	// 2. Call the Register RPC, passing the participant's local address.
	// 3. Handle the response and store the ID in the participant struct.
	return nil
}

// Deregister removes the participant from future question.
func (p *participant) Deregister() error {
	// TODO(student): Complete the method.
	// TODO(student): Remove comments below that are unnecessary when your implementation is complete.
	// TODO(student): Rephrase comments that help readability.
	// Complete the following steps:
	// 1. Using the connection object in Register and ID to call Deregister RPC
	return nil
}

// TODO(Student): Implement the Participant Service RPCs.
