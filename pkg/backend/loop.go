package backend

import (
	"context"

	"go-sqs-wrapper/pkg/workers"
)

// IOLoop structure
type IOLoop struct {
	QueueWorker *SQSQueue
	Worker      *workers.FutureWorker
	cancel      context.CancelFunc
}

// RunOnce starts one IO cycle
func (iol *IOLoop) RunOnce(ctx context.Context) {
	messages := iol.QueueWorker.FetchMessages(ctx)
	if len(messages) == 0 {
		return
	}
	processedMessages := iol.Worker.ProcessMessages(ctx, messages)
	if len(processedMessages) > 0 {
		iol.QueueWorker.DeleteMessages(ctx, processedMessages)
	}
}

// Run loops until ctx is cancelled (via the ctx or Stop), interrupting an
// in-flight fetch. It returns the ctx error that ended the loop.
func (iol *IOLoop) Run(ctx context.Context) error {
	ctx, iol.cancel = context.WithCancel(ctx)
	defer iol.cancel()
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		iol.RunOnce(ctx)
	}
}

// Stop cancels a running loop. Safe to call before Run (no-op) or multiple times.
func (iol *IOLoop) Stop() {
	if iol.cancel != nil {
		iol.cancel()
	}
}
