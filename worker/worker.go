package worker

import (
	"log"
	"os"
	"sync"

	"github.com/crowdmob/goamz/sqs"
)

var defaultStackName = "golang-sqsl-worker-example"

type HandlerFunc func(msg *sqs.Message) error

func (f HandlerFunc) HandleMessage(msg *sqs.Message) error {
	return f(msg)
}

type Handler interface {
	HandleMessage(msg *sqs.Message) error
}

func Start(queue *sqs.Queue, h Handler) {
	for {
		log.Println("worker: Start polling")
		poll(queue, h)
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

func handleMessage(h Handler, q *sqs.Queue, m *sqs.Message) error {
	var err error
	err = h.HandleMessage(m)
	if err != nil {
		return err
	}
	// delete
	_, err = q.DeleteMessage(m)
	if err != nil {
		return err
	}
	return nil
}

func poll(queue *sqs.Queue, h Handler) error {
	messages, err := queue.ReceiveMessage(10)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(messages.Messages))
	for i := range messages.Messages {
		// launch goroutine
		log.Println("worker: Spawn worker goroutine")

		go func(m *sqs.Message) {
			defer wg.Done()
			if err := handleMessage(h, queue, m); err != nil {
				log.Println(err)
			}
		}(&messages.Messages[i])
	}

	wg.Wait()
	return nil
}
