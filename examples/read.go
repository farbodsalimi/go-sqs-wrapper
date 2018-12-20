package examples

import (
	"../src/backend"
	"../src/util"
	"../src/workers"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"time"
)

func Read() {
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
	// Reading from SQS
	//
	fw := workers.FutureWorker{
		Handler: handler,
		TimeOut: 5,
	}

	//
	// Interrupting the reading process
	//
	time.AfterFunc(3*time.Second, fw.Stop)

	backend.IOLoop{
		QueueWorker: sqsWorker,
		Worker:      fw,
		StopSignal:  false,
	}.Run()
}

func handler(msg *sqs.Message) *sqs.Message {
	log.Println(msg)
	return msg
}
