package backend

import "go-sqs-wrapper/pkg/workers"

// IOLoop structure
type IOLoop struct {
	QueueWorker *SQSQueue
	Worker      *workers.FutureWorker
}

// RunOnce starts one IO cycle
func (iol *IOLoop) RunOnce() {
	messages := iol.QueueWorker.FetchMessages()
	if len(messages) > 0 {
		processedMessages := iol.Worker.ProcessMessages(messages)
		if len(processedMessages) > 0 {
			iol.QueueWorker.DeleteMessages(processedMessages)
		}
	}
}

// Run starts an infinite IO loop until the worker's context is cancelled
func (iol *IOLoop) Run() {
	for iol.Worker.IsRunning() {
		iol.RunOnce()
	}
}
