package examples

import (
	"../src/backend"
	"../src/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
)

func Publish() {
	util.LoadSettings()
	Region := os.Getenv("AWS_REGION")
	QueueURL := os.Getenv("SQS_QUEUE_URL")
	CredPath := os.Getenv("AWS_CRED_PATH")
	CredProfile := os.Getenv("AWS_CRED_PROFILE")

	sqsWorker := backend.SQSWorker{
		Region:      Region,
		QueueUrl:    QueueURL,
		CredPath:    CredPath,
		CredProfile: CredProfile,
	}
	sqsWorker.Init()

	//
	// Publishing to SQS
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

	sqsWorker.Publish("message 1", mav)
	sqsWorker.Publish("message 2", mav)
	sqsWorker.Publish("first message 3", mav)
	sqsWorker.Publish("first message 4", mav)
}
