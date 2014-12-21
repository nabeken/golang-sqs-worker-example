package worker

import (
	"log"
	"sync"

	"github.com/nabeken/aws-go-sqs/queue"
	"github.com/nabeken/aws-go-sqs/queue/option"
	"github.com/stripe/aws-go/gen/sqs"
)

var defaultStackName = "golang-sqsl-worker-example"

type HandlerFunc func(msg *sqs.Message) error

func (f HandlerFunc) HandleMessage(msg *sqs.Message) error {
	return f(msg)
}

type Handler interface {
	HandleMessage(msg *sqs.Message) error
}

func Start(q *queue.Queue, h Handler) {
	for {
		log.Println("worker: Start polling")
		messages, err := q.ReceiveMessage(option.MaxNumberOfMessages(10))
		if err != nil {
			log.Println(err)
			continue
		}
		if len(messages) > 0 {
			run(q, h, messages)
		}
	}
}

// poll launches goroutine per received message and wait for all message to be processed
func run(q *queue.Queue, h Handler, messages []sqs.Message) {
	numMessages := len(messages)
	log.Printf("worker: Received %d messages", numMessages)

	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *sqs.Message) {
			// launch goroutine
			log.Println("worker: Spawned worker goroutine")
			defer wg.Done()
			if err := handleMessage(q, m, h); err != nil {
				log.Println(err)
			}
		}(&messages[i])
	}

	wg.Wait()
}

func handleMessage(q *queue.Queue, m *sqs.Message, h Handler) error {
	var err error
	err = h.HandleMessage(m)
	if err != nil {
		return err
	}
	return q.DeleteMessage(m.ReceiptHandle)
}
