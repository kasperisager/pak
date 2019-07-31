package main

import (
	"os"

	"github.com/kasperisager/pak/cmd/pak/internal/build"
	"github.com/kasperisager/pak/pkg/cli"
)

func main() {
	app := cli.New("pak", "<command> [<arguments>]")

	app.AddCommand("build", "Build the thing!", build.Command)

	app.Run(os.Args[1:])
}
