package workers

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type futureMessage chan *sqs.Message

type FutureWorker struct {
	Handler     func(msg *sqs.Message) *sqs.Message
	TimeOut     int
	concurrency int
}

func (fw FutureWorker) ProcessMessages(messages []*sqs.Message) []*sqs.Message {
	var results []*sqs.Message
	for _, msg := range messages {
		results = append(results, <-processMessageAsync(fw.Handler, msg))
	}
	return results
}

func processMessageAsync(handler func(msg *sqs.Message) *sqs.Message, msg *sqs.Message) futureMessage {
	future := make(futureMessage)
	go func() { future <- handler(msg) }()
	return future
}