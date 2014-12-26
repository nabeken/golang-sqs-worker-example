package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/stripe/aws-go/aws"
	"github.com/stripe/aws-go/gen/sqs"

	"github.com/nabeken/golang-sqs-worker-example/worker"
)

var flagQueue = flag.String("q", "example", "specify a queue name")

func Print(msg *sqs.Message) error {
	fmt.Println(*msg.Body)
	return nil
}

func main() {
	flag.Parse()

	q, err := worker.NewSQSQueue(
		sqs.New(aws.DetectCreds("", "", ""), "ap-northeast-1", nil),
		*flagQueue,
	)
	if err != nil {
		log.Fatal(err)
	}

	worker.Start(q, worker.HandlerFunc(Print))
}
