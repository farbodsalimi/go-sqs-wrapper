package backend

import (
	"../workers"
)

type IOLoop struct {
	StopSignal  bool
	QueueWorker SQSQueue
	Worker      workers.FutureWorker
}

func (iol IOLoop) RunOnce() {
	messages := iol.QueueWorker.FetchMessages()
	if len(messages) > 0 {
		processedMessages := iol.Worker.ProcessMessages(messages)
		if len(processedMessages) > 0 {
			iol.QueueWorker.DeleteMessages(processedMessages)
		}
	}
}

func (iol IOLoop) Run() {
	for iol.StopSignal == false {
		iol.RunOnce()
	}
}
