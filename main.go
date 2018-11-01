package main

import (
	"fmt"
	"os"

	"./backend"
	"./util"
)

func main() {
	util.LoadSettings()
	Region := os.Getenv("AWS_REGION")
	QueueURL := os.Getenv("SQS_QUEUE_URL")
	CredPath := os.Getenv("AWS_CRED_PATH")
	CredProfile := os.Getenv("AWS_CRED_PROFILE")

	sqs := backend.SQSWorker{
		Region:      Region,
		QueueUrl:    QueueURL,
		CredPath:    CredPath,
		CredProfile: CredProfile,
	}

	sqs.Init()
	sqs.Publish("message 1", nil)
	sqs.Publish("message 2", nil)
	sqs.Publish("first message 3", nil)
	sqs.Publish("first message 4", nil)
	messages := sqs.FetchMessages()

	for _, value := range messages {
		fmt.Println(value)
	}
}
