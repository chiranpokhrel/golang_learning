package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"dat520/lab2/poll"
)

func main() {
	srvAddr := flag.String("srvAddr", "", "The address of the server")
	flag.Parse()

	pm, err := poll.NewPollMaster(*srvAddr)
	if err != nil {
		log.Fatalf("Unable to create PollMaster: %v", err)
	}
	log.Printf("Listening on: %v \n", pm.SrvAddress())

	questionId := 0
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Please write a question. Confirm question by pressing enter. To quit press enter:")
		question, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Unable to read: %v \n", err)
		}
		question = strings.Trim(question, "\n")
		if question == "" {
			log.Println("Closing service")
			break
		}
		options := []string{}
		fmt.Println("Please write options. Confirm option by pressing enter. To stop writing options, write an empty line.")
		for {
			option, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Unable to read: %v \n", err)
			}
			if option == "\n" {
				if len(options) < 2 {
					fmt.Println("You need to provide at least two options")
				} else {
					break
				}
			} else {
				option = strings.Trim(option, "\n")
				options = append(options, option)
			}
		}
		q := poll.Question{
			ID:       questionId,
			Question: question,
			Options:  options,
		}
		fmt.Printf("Sending question: %+v \n", q)
		questionId += 1
		resp, err := pm.HandleQuestion(&q)
		if err != nil {
			log.Printf("Unable to send question: %v \n", err)
		}

		if resp == nil {
			fmt.Println("Did not receive response from enough participants")
		} else if resp.Answer == poll.InconclusiveAnswer {
			fmt.Println("Received response: InconclusiveAnswer")
		} else {
			fmt.Printf("Received response: %+v \n", q.Options[resp.Answer])
		}
	}
}

type RandomQuestionHandler struct {
	Id int
}

func (rqh *RandomQuestionHandler) HandleQuestion(q *poll.Question) (*poll.Response, error) {
	answer := rand.Intn(len(q.Options))
	log.Printf("QH: %v Answering: %v to question: %v", rqh.Id, answer, q.ID)
	return &poll.Response{
		ID:     q.ID,
		Answer: answer,
	}, nil
}
