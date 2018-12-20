package main

import (
	"./src/cli"
	"bytes"
	"fmt"
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

	fmt.Printf("Loading handler: %v ...\n", cli.Opts.Handler)

	command := exec.Command("ls")

	// set var to get the output
	var out bytes.Buffer

	// set the output to our variable
	command.Stdout = &out
	err = command.Run()
	if err != nil {
		log.Println(err)
	}

	fmt.Println(out.String())
}
