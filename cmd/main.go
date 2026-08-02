package main

import (
	"os"

	"github.com/gllm-dev/molinets/cmd/cli"
)

func main() {
	if err := cli.NewCmdRoot().Execute(); err != nil {
		os.Exit(1)
	}
}
