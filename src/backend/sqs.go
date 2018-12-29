package backend

import (
	"../util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

var Client *sqs.SQS

// SQSQueue structure
type SQSQueue struct {
	QueueUrl            string
	Region              string
	MaxRetries          int
	MaxNumberOfMessages int64
	VisibilityTimeout   int64
	WaitTimeSeconds     int64
}

// Init, initializes the sqs client object
func (sq SQSQueue) Init() SQSQueue {
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

	queueUrl, err := util.GetEnvStr(util.GetEnvParams{Key: "SQS_QUEUE_URL"})
	util.CheckErr(err)
	sq.QueueUrl = queueUrl

	maxRetries, err := util.GetEnvInt(util.GetEnvParams{Key: "SQS_MAX_RETRIES"})
	util.CheckErr(err)
	sq.MaxRetries = maxRetries

	maxNumberOfMessages, err := util.GetEnvInt64(util.GetEnvParams{"MAX_NUMBER_OF_MESSAGES", "1"})
	util.CheckErr(err)
	sq.MaxNumberOfMessages = maxNumberOfMessages

	visibilityTimeout, err := util.GetEnvInt64(util.GetEnvParams{"VISIBILITY_TIMEOUT", "20"})
	util.CheckErr(err)
	sq.VisibilityTimeout = visibilityTimeout

	waitTimeSeconds, err := util.GetEnvInt64(util.GetEnvParams{"WAIT_TIME_SECONDS", "0"})
	util.CheckErr(err)
	sq.WaitTimeSeconds = waitTimeSeconds

	/*
	Create a new AWS session
	*/
	s, err := session.NewSession(&aws.Config{
		Region:     aws.String(sq.Region),
		MaxRetries: aws.Int(sq.MaxRetries),
	})
	if err != nil {
		log.Fatal("SQSQueue NewSession Error", err)
	}
	Client = sqs.New(s)

	return sq
}

// FetchMessages, fetches messages from sqs
func (sq SQSQueue) FetchMessages() []*sqs.Message {
	result, err := Client.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(sq.QueueUrl),
		MaxNumberOfMessages: aws.Int64(sq.MaxNumberOfMessages),
		VisibilityTimeout:   aws.Int64(sq.VisibilityTimeout), // 20 seconds
		WaitTimeSeconds:     aws.Int64(sq.WaitTimeSeconds),
	})

	if err != nil {
		log.Println("Fetch Messages Error", err)
	}

	if len(result.Messages) == 0 {
		log.Println("Received no messages")
	}

	return result.Messages
}

// DeleteMessages, deletes messages from sqs
func (sq SQSQueue) DeleteMessages(messages []*sqs.Message) {
	for _, message := range messages {
		resultDelete, err := Client.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(sq.QueueUrl),
			ReceiptHandle: message.ReceiptHandle,
		})

		if err != nil {
			log.Println("Delete Error", err)
			return
		}

		log.Println("Successfully deleted the message!", resultDelete)
	}
}

// Publish, publishes messages to sqs
func (sq SQSQueue) Publish(message string, mav map[string]*sqs.MessageAttributeValue) {

	var messageAttributes map[string]*sqs.MessageAttributeValue
	if messageAttributes = mav; messageAttributes == nil {
		messageAttributes = map[string]*sqs.MessageAttributeValue{}
	}

	result, err := Client.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:      aws.Int64(10),
		MessageAttributes: messageAttributes,
		MessageBody:       aws.String(message),
		QueueUrl:          aws.String(sq.QueueUrl),
	})

	if err != nil {
		log.Println("Error", err)
		return
	}

	log.Println("Successfully published to the queue!", *result.MessageId)
}
