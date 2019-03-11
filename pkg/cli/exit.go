package cli

import (
	"os"
	"sync"
)

var (
	exitCode = 0
	exitLock sync.Mutex
)

func Exit() {
	os.Exit(exitCode)
}

func ExitCode(code int) {
	exitLock.Lock()
	defer exitLock.Unlock()

	if exitCode < code {
		exitCode = code
	}
}
