package backend

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var Client *sqs.SQS

type SQSWorker struct {
	QueueUrl    string
	Region      string
	CredPath    string
	CredProfile string
}

// Init, initializes the sqs client object
func (sw SQSWorker) Init() {

	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(sw.Region),
		Credentials: credentials.NewSharedCredentials(sw.CredPath, sw.CredProfile),
		MaxRetries:  aws.Int(5),
	})

	if err != nil {
		fmt.Println("Init Error", err)
		return
	}

	// Create a SQS service client.
	Client = sqs.New(s)
}

// FetchMessages, fetches messages from sqs
func (sw SQSWorker) FetchMessages() []*sqs.Message {
	result, err := Client.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(sw.QueueUrl),
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(20), // 20 seconds
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		fmt.Println("Fetch Messages Error", err)
	}

	if len(result.Messages) == 0 {
		fmt.Println("Received no messages")
	}

	return result.Messages
}

// DeleteMessages, deletes messages from sqs
func (sw SQSWorker) DeleteMessages(messages []sqs.Message) {
	for _, message := range messages {
		resultDelete, err := Client.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(sw.QueueUrl),
			ReceiptHandle: message.ReceiptHandle,
		})

		if err != nil {
			fmt.Println("Delete Error", err)
			return
		}

		fmt.Println("Message Deleted", resultDelete)
	}
}

// Publish, publishes messages to sqs
func (sw SQSWorker) Publish(message string, mav map[string]*sqs.MessageAttributeValue) {

	var messageAttributes map[string]*sqs.MessageAttributeValue
	if messageAttributes = mav; mav == nil {
		messageAttributes = map[string]*sqs.MessageAttributeValue{}
	}

	result, err := Client.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:      aws.Int64(10),
		MessageAttributes: messageAttributes,
		MessageBody:       aws.String(message),
		QueueUrl:          aws.String(sw.QueueUrl),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success", *result.MessageId)
}
