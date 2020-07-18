package examples

import (
	"log"
	"time"

	"go-sqs-wrapper/pkg/backend"
	"go-sqs-wrapper/pkg/workers"

	"github.com/aws/aws-sdk-go/service/sqs"
)

// Read an example of reading messages from sqs
func Read() {
	//
	// Initialize your SQS queue.
	//
	sqsQueue := new(backend.SQSQueue).Init()

	//
	// Creating a future worker
	//
	fw := workers.FutureWorker{
		Handler: handler,
		TimeOut: 5,
	}

	//
	// Interrupting the reading process
	//
	time.AfterFunc(3*time.Second, fw.Stop)

	//
	// Starting an infinite io loop
	//
	backend.IOLoop{
		QueueWorker: sqsQueue,
		Worker:      fw,
		StopSignal:  false,
	}.Run()
}

func handler(msg *sqs.Message) *sqs.Message {
	log.Println(msg)
	return msg
}
