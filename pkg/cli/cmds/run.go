package cmds

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"

	"go-sqs-wrapper/pkg/backend"
	"go-sqs-wrapper/pkg/util"
	"go-sqs-wrapper/pkg/workers"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/subcommands"
)

// RunCmd struct
type RunCmd struct {
	handler string
}

// Name returns run subcommand's name
func (*RunCmd) Name() string { return "run" }

// Synopsis returns run subcommand's description
func (*RunCmd) Synopsis() string { return "start an IO Loop with the given message handler" }

// Usage returns run subcommand's usage
func (*RunCmd) Usage() string {
	return `run [-handler] <path to your handler>:
  Start IO Loop with the given message handler.
`
}

// SetFlags for run subcommands
func (rc *RunCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&rc.handler, "handler", "", "message handler")
}

// Execute run subcommand
func (rc *RunCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if rc.handler == "" || !util.IsCommandAvailable(rc.handler) {
		fmt.Println("Your handler is not valid!")
		return subcommands.ExitFailure
	}

	fmt.Printf("Selected message handler: %s\n\n", rc.handler)

	fw := workers.FutureWorker{
		Handler: func(msg *sqs.Message) *sqs.Message {
			// convert sqs message to string
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Fatal(err)
			}
			strMsg := string(jsonMsg)

			// pass the message to the handler and execute it
			command := exec.Command(rc.handler, strMsg)

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

	return subcommands.ExitSuccess
}
