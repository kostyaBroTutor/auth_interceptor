// Package process provides helper functions for managing the execution process.
package process

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForTermination waits the termination signal and allows to finish the program normally.
func WaitForTermination() {
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGTERM)
	<-stopSignal
}
