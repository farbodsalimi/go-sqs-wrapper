package workers

import (
	"log"
	"syscall"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type futureMessage chan *sqs.Message

// FutureWorker structure
type FutureWorker struct {
	Handler func(msg *sqs.Message) *sqs.Message
	TimeOut int
}

// ProcessMessages processes incoming messages async
func (fw FutureWorker) ProcessMessages(messages []*sqs.Message) []*sqs.Message {
	var results []*sqs.Message
	for _, msg := range messages {
		results = append(results, <-processMessageAsync(fw.Handler, msg))
	}
	return results
}

// Stop kills the worker
func (fw FutureWorker) Stop() {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		log.Panic(err)
	}
}

func processMessageAsync(handler func(msg *sqs.Message) *sqs.Message, msg *sqs.Message) futureMessage {
	future := make(futureMessage)
	go func() { future <- handler(msg) }()
	return future
}
