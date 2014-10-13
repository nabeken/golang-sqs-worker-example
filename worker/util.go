package worker

import (
	"os"

	"github.com/crowdmob/goamz/sqs"
)

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
