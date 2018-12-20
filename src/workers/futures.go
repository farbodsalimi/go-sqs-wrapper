package workers

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"syscall"
)

type futureMessage chan *sqs.Message

type FutureWorker struct {
	Handler func(msg *sqs.Message) *sqs.Message
	TimeOut int
}

func (fw FutureWorker) ProcessMessages(messages []*sqs.Message) []*sqs.Message {
	var results []*sqs.Message
	for _, msg := range messages {
		results = append(results, <-processMessageAsync(fw.Handler, msg))
	}
	return results
}

func (fw FutureWorker) Stop() {
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}

func processMessageAsync(handler func(msg *sqs.Message) *sqs.Message, msg *sqs.Message) futureMessage {
	future := make(futureMessage)
	go func() { future <- handler(msg) }()
	return future
}
