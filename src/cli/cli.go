package cli

// Opts structure which defines CLI operations
var Opts struct {
	Handler string `long:"handler" description:"The message handler that will be used by the workers."`
}
