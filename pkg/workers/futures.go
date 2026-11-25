package workers

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type futureMessage chan *types.Message

// FutureWorker structure
type FutureWorker struct {
	Handler    func(ctx context.Context, msg *types.Message) *types.Message
	TimeOut    int
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// ProcessMessages processes incoming messages concurrently using fan-out/fan-in.
func (fw *FutureWorker) ProcessMessages(messages []types.Message) []types.Message {
	ctx := fw.ctx
	if ctx == nil {
		ctx = context.Background()
	}

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

// Stop gracefully stops the worker
func (fw *FutureWorker) Stop() {
	if fw.cancelFunc != nil {
		fw.cancelFunc()
	}
}

// WithContext sets up the worker with a cancellable context
func (fw *FutureWorker) WithContext(ctx context.Context) {
	fw.ctx, fw.cancelFunc = context.WithCancel(ctx)
}

// IsRunning returns true if the worker context is not cancelled
func (fw *FutureWorker) IsRunning() bool {
	if fw.ctx == nil {
		return true // no context means always running for backward compatibility
	}
	select {
	case <-fw.ctx.Done():
		return false
	default:
		return true
	}
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
