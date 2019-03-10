package cli

import (
	"os"
)

var exitCode = 0

func exit() {
	os.Exit(exitCode)
}
