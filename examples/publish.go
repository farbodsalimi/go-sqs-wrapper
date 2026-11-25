package examples

import (
	"go-sqs-wrapper/pkg/backend"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Publish an example of publishing messages to sqs
func Publish() {
	sqsQueue := new(backend.SQSQueue).Init()

	//
	// Publishing to SQSQueue
	//
	mav := map[string]types.MessageAttributeValue{
		"Title": {
			DataType:    new("String"),
			StringValue: new("The Whistler"),
		},
		"Author": {
			DataType:    new("String"),
			StringValue: new("Farbod Salimi"),
		},
		"WeeksOn": {
			DataType:    new("Number"),
			StringValue: new("6"),
		},
	}

	sqsQueue.Publish("message 1", mav)
	sqsQueue.Publish("message 2", mav)
	sqsQueue.Publish("message 3", mav)
	sqsQueue.Publish("message 4", mav)
}
