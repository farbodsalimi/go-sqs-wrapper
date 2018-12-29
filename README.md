# golang-sqs-wrapper

## CLI

### Help

```bash
./bin/sqs_wrapper --help
```

```bash
Usage:
  sqs_wrapper [OPTIONS]

Application Options:
      --handler= The message handler that will be used by the workers.

Help Options:
  -h, --help     Show this help message

panic: Usage:
  sqs_wrapper [OPTIONS]

Application Options:
      --handler= The message handler that will be used by the workers.

Help Options:
  -h, --help     Show this help message
```

### Run the IO loop with your handler

```
./bin/sqs_wrapper --handler={path to your handler or your command}
```

Examples:

```
./bin/sqs_wrapper --handler=echo
```

```
./bin/sqs_wrapper --handler=~/path/to/your/handler/main
```

## Code Examples

```go
sqsQueue := new(backend.SQSQueue).Init()

//
// Publishing to SQS
//
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

//
// Reading from SQS
//
fw := workers.FutureWorker{
    Handler: func (msg *sqs.Message) *sqs.Message {
                log.Println(msg)
                return msg
             },
    TimeOut: 5,
}

//
// Interrupting the reading process
//
time.AfterFunc(3*time.Second, fw.Stop)

backend.IOLoop{
    QueueWorker: sqsQueue,
    Worker:      fw,
    StopSignal:  false,
}.Run()
```
