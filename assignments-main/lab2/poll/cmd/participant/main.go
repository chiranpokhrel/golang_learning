package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"dat520/lab2/poll"
)

func main() {
	srvAddr := flag.String("srvAddr", "", "The address of the server")
	flag.Parse()

	answerChan := make(chan int)
	qh := NewCmdQuestionHandler(answerChan)
	pp, err := poll.NewParticipant(qh, "")
	if err != nil {
		log.Fatalf("Unable to create participant: %v", err)
	}

	err = pp.Register(*srvAddr)
	if err != nil {
		log.Fatalf("Unable to register participant: %v", err)
	}

	log.Println("Successfully registered participant")
	log.Println("Participants waiting for questions")

	fmt.Println("To close the participant write 'exit' at any point.")

	cmdChan := make(chan string)
	go ReadFromCmd(cmdChan)
	for {
		input, ok := <-cmdChan
		if !ok {
			break
		}
		switch input {
		case "exit":
			log.Println("Closing participant")
			close(cmdChan)
			close(answerChan)
			err := pp.Deregister()
			if err != nil {
				log.Fatalf("Unable to deregister participant: %v", err)
			}
			log.Println("Deregistered participant")
		default:
			answer, err := strconv.Atoi(input)
			if err != nil {
				fmt.Printf("Unable to parse response: %v\n", err)
				continue
			}
			fmt.Println("Sending value", answer)
			answerChan <- answer
		}
	}
}

func ReadFromCmd(outChan chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Unable to read response: %v", err)
			continue
		}
		response = strings.Trim(response, "\n")
		outChan <- response
	}
}

type CmdQuestionHandler struct {
	answerChan chan int
	qChan      chan *poll.Question
}

func NewCmdQuestionHandler(answerChan chan int) *CmdQuestionHandler {
	qChan := make(chan *poll.Question)
	cmd := &CmdQuestionHandler{
		answerChan: answerChan,
		qChan:      qChan,
	}
	go cmd.writeQuestion(qChan)
	return cmd
}

func (cmd *CmdQuestionHandler) HandleQuestion(q *poll.Question) (*poll.Response, error) {
	cmd.qChan <- q
	answer := poll.InconclusiveAnswer
	var ok bool
	select {
	case answer, ok = <-cmd.answerChan:
		if !ok {
			return nil, fmt.Errorf("closing participant")
		}
	case <-time.After(poll.QuestionTimeout):
		fmt.Println("Question timed out.")
	}
	return &poll.Response{
		ID:     q.ID,
		Answer: answer,
	}, nil
}

func (cmd *CmdQuestionHandler) writeQuestion(qChan chan *poll.Question) {
	for {
		q, ok := <-qChan
		if !ok {
			break
		}
		fmt.Printf("Question: %v \n", q.Question)
		for i, option := range q.Options {
			fmt.Printf("\tOption %v: %v\n", i, option)
		}
	}
}
