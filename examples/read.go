package examples

import (
	"../src/backend"
	"../src/workers"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"time"
)

func Read() {
	sqsQueue := new(backend.SQSQueue).Init()

	//
	// Reading from SQSQueue
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
		QueueWorker: sqsQueue,
		Worker:      fw,
		StopSignal:  false,
	}.Run()
}

func handler(msg *sqs.Message) *sqs.Message {
	log.Println(msg)
	return msg
}
