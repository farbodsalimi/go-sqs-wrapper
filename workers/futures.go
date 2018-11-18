package workers

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

type futureMessage chan string

type FutureWorker struct {
	Handler     func(msg *sqs.Message) string
	TimeOut     int
	concurrency int
}

func (fw FutureWorker) ProcessMessages(messages []*sqs.Message) []string {
	var results []string
	for _, msg := range messages {
		results = append(results, <-processMessageAsync(fw.Handler, msg))
	}
	return results
}

func processMessageAsync(handler func(msg *sqs.Message) string, msg *sqs.Message) futureMessage {
	future := make(futureMessage)
	go func() { future <- handler(msg) }()
	return future
}
