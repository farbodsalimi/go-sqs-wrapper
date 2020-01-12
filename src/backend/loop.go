package backend

import "go-sqs-wrapper/src/workers"

// IOLoop structure
type IOLoop struct {
	StopSignal  bool
	QueueWorker SQSQueue
	Worker      workers.FutureWorker
}

// RunOnce starts one IO cycle
func (iol IOLoop) RunOnce() {
	messages := iol.QueueWorker.FetchMessages()
	if len(messages) > 0 {
		processedMessages := iol.Worker.ProcessMessages(messages)
		if len(processedMessages) > 0 {
			iol.QueueWorker.DeleteMessages(processedMessages)
		}
	}
}

// Run starts an infinite IO loop
func (iol IOLoop) Run() {
	for !iol.StopSignal {
		iol.RunOnce()
	}
}
