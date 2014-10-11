package worker

import (
	"log"
	"os"

	"github.com/crowdmob/goamz/sqs"
)

var defaultStackName = "golang-sqsl-worker-example"

type HandlerFunc func(msg *sqs.Message) bool

func (f HandlerFunc) HandleMessage(msg *sqs.Message) bool {
	return f(msg)
}

type Handler interface {
	HandleMessage(msg *sqs.Message) bool
}

func Start(queue *sqs.Queue, h Handler) error {
	for {
		log.Println("worker: Start polling")
		if err := poll(queue, h); err != nil {
			return err
		}
	}
}

func Getenv(name, defaultVal string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal
	}
	return val
}

func NewSQSQueue(s *sqs.SQS, name string) (*sqs.Queue, error) {
	stackName := Getenv("AWS_STACK_NAME", defaultStackName)
	if stackName != "" {
		stackName += "-"
	}
	return s.GetQueue(stackName + name)
}

func poll(queue *sqs.Queue, h Handler) error {
	messages, err := queue.ReceiveMessage(10)
	if err != nil {
		return err
	}
	for i := range messages.Messages {
		log.Println("worker: Process message")
		if !h.HandleMessage(&messages.Messages[i]) {
			// will retry
			continue
		}
		// delete
		if _, err := queue.DeleteMessage(&messages.Messages[i]); err != nil {
			return err
		}
	}
	return nil
}
