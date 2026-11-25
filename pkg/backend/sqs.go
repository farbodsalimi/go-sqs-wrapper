package backend

import (
	"context"
	"log"

	"go-sqs-wrapper/pkg/util"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// SQSQueue structure
type SQSQueue struct {
	QueueURL            string
	Region              string
	MaxNumberOfMessages int32
	VisibilityTimeout   int32
	WaitTimeSeconds     int32
	DelaySeconds        int32
	client              *sqs.Client
}

// NewSQSQueue creates an SQSQueue with the given client and configuration.
func NewSQSQueue(client *sqs.Client, queueURL string) *SQSQueue {
	return &SQSQueue{
		QueueURL:            queueURL,
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   20,
		WaitTimeSeconds:     0,
		DelaySeconds:        0,
		client:              client,
	}
}

// Init initializes the sqs client object
func (sq *SQSQueue) Init() *SQSQueue {
	/*
		Load all the required config from the .env file.
		Consider that AWS SDK will automatically use the
		AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY from
		the environment therefore we don't have to manually
		pass them along.
	*/
	util.LoadDotEnv()

	region, err := util.GetEnvStr(util.GetEnvParams{Key: "AWS_REGION"})
	sq.Region = region
	util.CheckErr(err)

	queueURL, err := util.GetEnvStr(util.GetEnvParams{Key: "SQS_QUEUE_URL"})
	util.CheckErr(err)
	sq.QueueURL = queueURL

	maxNumberOfMessages, err := util.GetEnvInt(util.GetEnvParams{Key: "MAX_NUMBER_OF_MESSAGES", DefaultValue: "1"})
	util.CheckErr(err)
	sq.MaxNumberOfMessages = int32(maxNumberOfMessages)

	visibilityTimeout, err := util.GetEnvInt(util.GetEnvParams{Key: "VISIBILITY_TIMEOUT", DefaultValue: "20"})
	util.CheckErr(err)
	sq.VisibilityTimeout = int32(visibilityTimeout)

	waitTimeSeconds, err := util.GetEnvInt(util.GetEnvParams{Key: "WAIT_TIME_SECONDS", DefaultValue: "0"})
	util.CheckErr(err)
	sq.WaitTimeSeconds = int32(waitTimeSeconds)

	delaySeconds, err := util.GetEnvInt(util.GetEnvParams{Key: "DELAY_SECONDS", DefaultValue: "10"})
	util.CheckErr(err)
	sq.DelaySeconds = int32(delaySeconds)

	/*
		Create a new AWS config and SQS client
	*/
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(sq.Region),
	)
	if err != nil {
		log.Fatal("SQS LoadDefaultConfig Error", err)
	}
	sq.client = sqs.NewFromConfig(cfg)

	return sq
}

// FetchMessages fetches messages from sqs
func (sq *SQSQueue) FetchMessages(ctx context.Context) []types.Message {
	result, err := sq.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{"All"},
		MessageSystemAttributeNames: []types.MessageSystemAttributeName{
			types.MessageSystemAttributeNameSentTimestamp,
		},
		QueueUrl:            &sq.QueueURL,
		MaxNumberOfMessages: sq.MaxNumberOfMessages,
		VisibilityTimeout:   sq.VisibilityTimeout,
		WaitTimeSeconds:     sq.WaitTimeSeconds,
	})

	if err != nil {
		log.Println("Fetch Messages Error", err)
		return nil
	}

	if len(result.Messages) == 0 {
		log.Println("Received no messages")
	}

	return result.Messages
}

// DeleteMessages deletes messages from sqs
func (sq *SQSQueue) DeleteMessages(ctx context.Context, messages []types.Message) {
	for _, message := range messages {
		_, err := sq.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &sq.QueueURL,
			ReceiptHandle: message.ReceiptHandle,
		})

		if err != nil {
			log.Println("Delete Error", err)
			continue
		}

		log.Println("Successfully deleted the message!")
	}
}

// Publish publishes messages to sqs
func (sq *SQSQueue) Publish(message string, mav map[string]types.MessageAttributeValue) {
	var messageAttributes map[string]types.MessageAttributeValue
	if messageAttributes = mav; messageAttributes == nil {
		messageAttributes = map[string]types.MessageAttributeValue{}
	}

	result, err := sq.client.SendMessage(context.Background(), &sqs.SendMessageInput{
		DelaySeconds:      sq.DelaySeconds,
		MessageAttributes: messageAttributes,
		MessageBody:       &message,
		QueueUrl:          &sq.QueueURL,
	})

	if err != nil {
		log.Println("Error", err)
		return
	}

	log.Println("Successfully published to the queue!", *result.MessageId)
}
