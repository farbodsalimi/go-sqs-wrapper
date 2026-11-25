//go:build integration

package integration

import (
	"context"
	"testing"

	"go-sqs-wrapper/pkg/backend"
	"go-sqs-wrapper/pkg/workers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const localstackEndpoint = "http://localhost:4566"

func newTestClient(t *testing.T) *sqs.Client {
	t.Helper()
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	require.NoError(t, err)

	return sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(localstackEndpoint)
	})
}

func createTestQueue(t *testing.T, client *sqs.Client, name string) string {
	t.Helper()
	out, err := client.CreateQueue(context.Background(), &sqs.CreateQueueInput{
		QueueName: aws.String(name),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		client.DeleteQueue(context.Background(), &sqs.DeleteQueueInput{
			QueueUrl: out.QueueUrl,
		})
	})
	return *out.QueueUrl
}

func TestPublishAndFetchMessages(t *testing.T) {
	client := newTestClient(t)
	queueURL := createTestQueue(t, client, "test-publish-fetch")

	sq := backend.NewSQSQueue(client, queueURL)
	sq.MaxNumberOfMessages = 10
	sq.WaitTimeSeconds = 1

	// Publish
	sq.Publish("hello world", nil)

	// Fetch
	messages := sq.FetchMessages()
	require.Len(t, messages, 1)
	assert.Equal(t, "hello world", *messages[0].Body)
}

func TestPublishWithAttributes(t *testing.T) {
	client := newTestClient(t)
	queueURL := createTestQueue(t, client, "test-publish-attrs")

	sq := backend.NewSQSQueue(client, queueURL)
	sq.WaitTimeSeconds = 1

	attrs := map[string]types.MessageAttributeValue{
		"Title": {
			DataType:    aws.String("String"),
			StringValue: aws.String("Test Title"),
		},
	}
	sq.Publish("with attrs", attrs)

	messages := sq.FetchMessages()
	require.Len(t, messages, 1)
	assert.Equal(t, "with attrs", *messages[0].Body)
	assert.Contains(t, messages[0].MessageAttributes, "Title")
	assert.Equal(t, "Test Title", *messages[0].MessageAttributes["Title"].StringValue)
}

func TestDeleteMessages(t *testing.T) {
	client := newTestClient(t)
	queueURL := createTestQueue(t, client, "test-delete")

	sq := backend.NewSQSQueue(client, queueURL)
	sq.WaitTimeSeconds = 1

	sq.Publish("to be deleted", nil)

	messages := sq.FetchMessages()
	require.Len(t, messages, 1)

	sq.DeleteMessages(messages)

	// Queue should be empty now
	remaining := sq.FetchMessages()
	assert.Empty(t, remaining)
}

func TestFetchEmptyQueue(t *testing.T) {
	client := newTestClient(t)
	queueURL := createTestQueue(t, client, "test-empty")

	sq := backend.NewSQSQueue(client, queueURL)
	sq.WaitTimeSeconds = 0

	messages := sq.FetchMessages()
	assert.Empty(t, messages)
}

func TestPublishMultipleMessages(t *testing.T) {
	client := newTestClient(t)
	queueURL := createTestQueue(t, client, "test-multi")

	sq := backend.NewSQSQueue(client, queueURL)
	sq.MaxNumberOfMessages = 10
	sq.WaitTimeSeconds = 1

	for i := 0; i < 3; i++ {
		sq.Publish("msg", nil)
	}

	var all []types.Message
	for range 3 {
		msgs := sq.FetchMessages()
		all = append(all, msgs...)
		if len(all) >= 3 {
			break
		}
	}
	assert.Len(t, all, 3)
}

func TestIOLoopProcessesMessages(t *testing.T) {
	client := newTestClient(t)
	queueURL := createTestQueue(t, client, "test-ioloop")

	sq := backend.NewSQSQueue(client, queueURL)
	sq.WaitTimeSeconds = 1

	sq.Publish("ioloop test", nil)

	var processed []string
	fw := &workers.FutureWorker{
		Handler: func(_ context.Context, msg *types.Message) *types.Message {
			processed = append(processed, *msg.Body)
			return msg
		},
		TimeOut: 5,
	}
	fw.WithContext(context.Background())

	// Run a single cycle
	iol := &backend.IOLoop{
		QueueWorker: sq,
		Worker:      fw,
	}
	iol.RunOnce()

	assert.Equal(t, []string{"ioloop test"}, processed)

	// Message should be deleted after processing
	remaining := sq.FetchMessages()
	assert.Empty(t, remaining)
}

func TestFutureWorkerConcurrency(t *testing.T) {
	client := newTestClient(t)
	queueURL := createTestQueue(t, client, "test-concurrency")

	sq := backend.NewSQSQueue(client, queueURL)
	sq.MaxNumberOfMessages = 10
	sq.WaitTimeSeconds = 1

	for i := 0; i < 5; i++ {
		sq.Publish("concurrent", nil)
	}

	fw := &workers.FutureWorker{
		Handler: func(_ context.Context, msg *types.Message) *types.Message {
			return msg
		},
		TimeOut: 5,
	}
	fw.WithContext(context.Background())

	// Fetch all messages (may need multiple fetches)
	var allMsgs []types.Message
	for range 5 {
		msgs := sq.FetchMessages()
		allMsgs = append(allMsgs, msgs...)
		if len(allMsgs) >= 5 {
			break
		}
	}

	// Process them all concurrently
	results := fw.ProcessMessages(allMsgs)
	assert.Len(t, results, len(allMsgs))
}

func TestFutureWorkerTimeout(t *testing.T) {
	client := newTestClient(t)
	queueURL := createTestQueue(t, client, "test-timeout")

	sq := backend.NewSQSQueue(client, queueURL)
	sq.WaitTimeSeconds = 1

	sq.Publish("timeout test", nil)

	fw := &workers.FutureWorker{
		Handler: func(ctx context.Context, msg *types.Message) *types.Message {
			// Block until context is cancelled (simulating a hung handler)
			<-ctx.Done()
			return nil
		},
		TimeOut: 1, // 1 second timeout
	}
	fw.WithContext(context.Background())

	messages := sq.FetchMessages()
	require.Len(t, messages, 1)

	// Should return nil results because handler times out
	results := fw.ProcessMessages(messages)
	assert.Empty(t, results)
}
