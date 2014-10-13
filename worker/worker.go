package worker

import (
	"log"
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

func Start(q *sqs.Queue, h Handler) {
	for {
		log.Println("worker: Start polling")
		resp, err := q.ReceiveMessage(10)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(resp.Messages) > 0 {
			run(q, h, resp)
		}
	}
}

// poll launches goroutine per received message and wait for all message to be processed
func run(q *sqs.Queue, h Handler, resp *sqs.ReceiveMessageResponse) {
	numMessages := len(resp.Messages)
	log.Printf("worker: Received %d messages", numMessages)

	var wg sync.WaitGroup
	wg.Add(numMessages)
	for i := range resp.Messages {
		go func(m *sqs.Message) {
			// launch goroutine
			log.Println("worker: Spawned worker goroutine")
			defer wg.Done()
			if err := handleMessage(q, m, h); err != nil {
				log.Println(err)
			}
		}(&resp.Messages[i])
	}

	wg.Wait()
}

func handleMessage(q *sqs.Queue, m *sqs.Message, h Handler) error {
	var err error
	err = h.HandleMessage(m)
	if err != nil {
		return err
	}
	// delete
	_, err = q.DeleteMessage(m)
	return err
}
