package main

import (
	"./src/backend"
	"./src/cli"
	"./src/util"
	"./src/workers"
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"os/exec"
)

func main() {
	_, err := flags.ParseArgs(&cli.Opts, os.Args)
	if err != nil {
		panic(err)
	}
	
	util.LoadSettings()
	Region := os.Getenv("AWS_REGION")
	QueueURL := os.Getenv("SQS_QUEUE_URL")
	CredPath := os.Getenv("AWS_CRED_PATH")
	CredProfile := os.Getenv("AWS_CRED_PROFILE")

	sqsWorker := backend.SQSWorker{
		Region:      Region,
		QueueUrl:    QueueURL,
		CredPath:    CredPath,
		CredProfile: CredProfile,
	}
	sqsWorker.Init()

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

	backend.IOLoop{
		QueueWorker: sqsWorker,
		Worker:      fw,
		StopSignal:  false,
	}.Run()
}
