package workers

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type futureMessage chan *types.Message

// FutureWorker structure
type FutureWorker struct {
	Handler func(ctx context.Context, msg *types.Message) *types.Message
	TimeOut int
}

// ProcessMessages processes incoming messages concurrently using fan-out/fan-in.
// A cancelled ctx aborts in-flight handlers (nil result, dropped from output).
func (fw *FutureWorker) ProcessMessages(ctx context.Context, messages []types.Message) []types.Message {
	// Fan-out: launch all goroutines first
	futures := make([]futureMessage, len(messages))
	for i := range messages {
		futures[i] = fw.processMessageAsync(ctx, &messages[i])
	}

	// Fan-in: collect results
	var results []types.Message
	for _, f := range futures {
		if msg := <-f; msg != nil {
			results = append(results, *msg)
		}
	}
	return results
}

func (fw *FutureWorker) processMessageAsync(ctx context.Context, msg *types.Message) futureMessage {
	// Buffered channel prevents goroutine leak if result is never read
	future := make(futureMessage, 1)

	timeout := time.Duration(fw.TimeOut) * time.Second
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	go func() {
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		done := make(chan *types.Message, 1)
		go func() {
			done <- fw.Handler(timeoutCtx, msg)
		}()

		select {
		case result := <-done:
			future <- result
		case <-timeoutCtx.Done():
			future <- nil
		}
	}()

	return future
}
