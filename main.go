package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"./src/backend"
	"./src/cli"
	"./src/workers"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jessevdk/go-flags"
)

func main() {
	_, err := flags.ParseArgs(&cli.Opts, os.Args)
	if err != nil {
		panic(err)
	}

	fw := workers.FutureWorker{
		Handler: func(msg *sqs.Message) *sqs.Message {
			// convert sqs message to string
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Fatal(err)
			}
			strMsg := string(jsonMsg)

			// pass the message to the handler and execute it
			command := exec.Command(cli.Opts.Handler, strMsg)

			// set var to get the output
			var out bytes.Buffer

			// set the output to our variable
			command.Stdout = &out
			err = command.Run()
			if err != nil {
				log.Println(err)
			}

			log.Println(out.String())
			return msg
		},
		TimeOut: 5,
	}

	sqsQueue := new(backend.SQSQueue).Init()

	backend.IOLoop{
		QueueWorker: sqsQueue,
		Worker:      fw,
		StopSignal:  false,
	}.Run()
}
