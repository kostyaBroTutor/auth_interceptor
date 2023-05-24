package process

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

// Terminate its application gratefully.
func Terminate() {
	pid := os.Getpid()

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Println(
			"failed to find process ID for terminate myself, error: " + err.Error(),
		)

		os.Exit(1)
	}

	if err = proc.Signal(syscall.SIGTERM); err != nil {
		log.Println(fmt.Sprintf(
			"failed to send SIGTERM signal to myself, PID: %d, error: %s",
			pid, err.Error(),
		))

		os.Exit(1)
	}
}
