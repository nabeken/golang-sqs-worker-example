package main

import (
	"fmt"
	"log"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/sqs"

	"github.com/nabeken/golang-sqs-worker-example/worker"
)

func Print(msg *sqs.Message) bool {
	fmt.Println(msg.Body)
	return true
}

func main() {
	auth, err := aws.GetAuth("", "", "", time.Now())
	if err != nil {
		log.Fatal(err)
	}

	s := sqs.New(auth, aws.APNortheast)
	queue, err := worker.NewSQSQueue(s, "example")
	if err != nil {
		log.Fatal(err)
	}

	worker.Start(queue, worker.HandlerFunc(Print))
}
