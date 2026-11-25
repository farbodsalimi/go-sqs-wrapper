package examples

import (
	"context"
	"log"
	"time"

	"go-sqs-wrapper/pkg/backend"
	"go-sqs-wrapper/pkg/workers"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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
	fw := &workers.FutureWorker{
		Handler: handler,
		TimeOut: 5,
	}
	fw.WithContext(context.Background())

	//
	// Interrupting the reading process
	//
	time.AfterFunc(3*time.Second, fw.Stop)

	//
	// Starting an infinite io loop
	//
	iol := &backend.IOLoop{
		QueueWorker: sqsQueue,
		Worker:      fw,
	}
	iol.Run()
}

func handler(_ context.Context, msg *types.Message) *types.Message {
	log.Println(msg)
	return msg
}
