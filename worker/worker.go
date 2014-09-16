package worker

import (
	"github.com/crowdmob/goamz/sqs"
	"log"
)

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
