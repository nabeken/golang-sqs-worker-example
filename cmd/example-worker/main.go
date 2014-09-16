package main

import (
	"flag"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/sqs"
	"log"
	"time"

	"github.com/nabeken/golang-sqs-worker-example/worker"
)

func Print(msg *sqs.Message) bool {
	fmt.Println(msg.Body)
	return true
}

func main() {
	queueName := flag.String("n", "", "Specify a queue name")
	flag.Parse()

	if *queueName == "" {
		log.Fatal("Queue name must not be string")
	}

	auth, err := aws.GetAuth("", "", "", time.Now())
	if err != nil {
		log.Fatal(err)
	}

	s := sqs.New(auth, aws.APNortheast)
	queue, err := s.GetQueue(*queueName)

	if err != nil {
		log.Fatal(err)
	}

	worker.Start(queue, worker.HandlerFunc(Print))
}
