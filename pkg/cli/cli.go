package cli

import (
	"context"
	"flag"
	"go-sqs-wrapper/pkg/cli/cmds"
	"os"

	"github.com/google/subcommands"
)

// Parse starts CLI
func Parse() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cmds.RunCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
