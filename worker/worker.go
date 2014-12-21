package worker

import (
	"log"
	"sync"

	"github.com/nabeken/aws-go-sqs/sqs"
	"github.com/nabeken/aws-go-sqs/sqs/option"

	gsqs "github.com/stripe/aws-go/gen/sqs"
)

var defaultStackName = "golang-sqsl-worker-example"

type HandlerFunc func(msg *gsqs.Message) error

func (f HandlerFunc) HandleMessage(msg *gsqs.Message) error {
	return f(msg)
}

type Handler interface {
	HandleMessage(msg *gsqs.Message) error
}

func Start(q *sqs.Queue, h Handler) {
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
func run(q *sqs.Queue, h Handler, messages []gsqs.Message) {
	numMessages := len(messages)
	log.Printf("worker: Received %d messages", numMessages)

	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range messages {
		go func(m *gsqs.Message) {
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

func handleMessage(q *sqs.Queue, m *gsqs.Message, h Handler) error {
	var err error
	err = h.HandleMessage(m)
	if err != nil {
		return err
	}
	return q.DeleteMessage(m.ReceiptHandle)
}
