# golang-sqs-wrapper

## Environment Variables

The first step is to define all of your environment variables in `.env`. You can find and example of `.env` file in the root folder.

| Name                   | Type   | Default |
| ---------------------- | ------ | ------- |
| AWS_ACCESS_KEY_ID      | string |         |
| AWS_SECRET_ACCESS_KEY  | string |         |
| AWS_REGION             | string |         |
| SQS_QUEUE_URL          | string |         |
| SQS_MAX_RETRIES        | string |         |
| MAX_NUMBER_OF_MESSAGES | number | 1       |
| VISIBILITY_TIMEOUT     | number | 20      |
| WAIT_TIME_SECONDS      | number | 0       |

## Usage

### - Initializing your sqs queue

```go
sqsQueue := new(backend.SQSQueue).Init()
```

### - Publishing to SQS

```go
mav := map[string]*sqs.MessageAttributeValue{
    "Title": {
        DataType:    aws.String("String"),
        StringValue: aws.String("The Whistler"),
    },
    "Author": {
        DataType:    aws.String("String"),
        StringValue: aws.String("Farbod Salimi"),
    },
    "WeeksOn": {
        DataType:    aws.String("Number"),
        StringValue: aws.String("6"),
    },
}
sqsQueue.Publish("message 1", mav)
```

### - Reading from SQS via CLI

```
./bin/sqs_wrapper run --handler=echo
```

```
./bin/sqs_wrapper run --handler=~/path/to/your/handler/main
```

```bash
./bin/sqs_wrapper --help

Usage: sqs_wrapper <flags> <subcommand> <subcommand args>

Subcommands:
        commands         list all command names
        flags            describe all known top-level flags
        help             describe subcommands and their syntax
        run              start an IO Loop with the given message handler


Use "sqs_wrapper flags" for a list of top-level flags
```

### - Reading from SQS programatically

```go
// Creating a future worker
fw := workers.FutureWorker{
    Handler: func (msg *sqs.Message) *sqs.Message {
                log.Println(msg)
                return msg
             },
    TimeOut: 5,
}

// Starting an infinite io loop
backend.IOLoop{
    QueueWorker: sqsQueue,
    Worker:      fw,
    StopSignal:  false,
}.Run()
```

Note: you can interrupt your IO loop by `fw.Stop()`, for example:

```bash
time.AfterFunc(3*time.Second, fw.Stop)`
```
