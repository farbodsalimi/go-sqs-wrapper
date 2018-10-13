package main

import (
	"fmt"

	"./backend"
)

const (
	QueueUrl    = "https://sqs.us-east-1.amazonaws.com/******/my-queue"
	Region      = "us-east-1"
	CredPath    = "/Users/username/.aws/credentials"
	CredProfile = "default"
)

func main() {
	sqs := backend.SQSWorker{
		Region:      Region,
		QueueUrl:    QueueUrl,
		CredPath:    CredPath,
		CredProfile: CredProfile,
	}
	sqs.Init()
	sqs.Publish("message 1")
	sqs.Publish("message 2")
	sqs.Publish("first message 3")
	sqs.Publish("first message 4")
	messages := sqs.FetchMessages()

	for _, value := range messages {
		fmt.Println(value)
	}
}
