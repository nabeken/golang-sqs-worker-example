package worker

import (
	"os"

	"github.com/nabeken/aws-go-sqs/sqs"

	gsqs "github.com/stripe/aws-go/gen/sqs"
)

func Getenv(name, defaultVal string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal
	}
	return val
}

func NewSQSQueue(s *gsqs.SQS, name string) (*sqs.Queue, error) {
	stackName := Getenv("AWS_STACK_NAME", defaultStackName)
	if stackName != "" {
		stackName += "-"
	}
	return sqs.New(s, stackName+name)
}
