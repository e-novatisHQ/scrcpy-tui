package main

import (
	"os"

	"github.com/e-novatisHQ/scrcpy-tui/internal/cli"
)

var version = "dev"

func main() {
	os.Exit(cli.Run(os.Args[1:], version, cli.Streams{In: os.Stdin, Out: os.Stdout, Err: os.Stderr}))
}
