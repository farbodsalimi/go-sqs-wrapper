package examples

import (
	"go-sqs-wrapper/pkg/backend"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Publish an example of publishing messages to sqs
func Publish() {
	sqsQueue := new(backend.SQSQueue).Init()

	//
	// Publishing to SQSQueue
	//
	mav := map[string]*sqs.MessageAttributeValue{
		"Title": {
			DataType:    aws.String("String"),
			StringValue: aws.String("The Whistler"),
		},
		"Author": {
			DataType:    aws.String("String"),
			StringValue: aws.String("Farbod Salimi"),
		},
		"WeeksOn": {
			DataType:    aws.String("Number"),
			StringValue: aws.String("6"),
		},
	}

	sqsQueue.Publish("message 1", mav)
	sqsQueue.Publish("message 2", mav)
	sqsQueue.Publish("message 3", mav)
	sqsQueue.Publish("message 4", mav)
}
