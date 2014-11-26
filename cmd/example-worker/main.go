package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/sqs"

	"github.com/nabeken/golang-sqs-worker-example/worker"
)

var flagQueue = flag.String("q", "example", "specify a queue name")

func Print(msg *sqs.Message) error {
	fmt.Println(msg.Body)
	return nil
}

func main() {
	flag.Parse()

	auth, err := aws.GetAuth("", "", "", time.Now())
	if err != nil {
		log.Fatal(err)
	}

	s := sqs.New(auth, aws.APNortheast)
	queue, err := worker.NewSQSQueue(s, *flagQueue)
	if err != nil {
		log.Fatal(err)
	}

	worker.Start(queue, worker.HandlerFunc(Print))
}
