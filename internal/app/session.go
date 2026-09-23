package app

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

type Interrupted struct{ Signal os.Signal }

func (e *Interrupted) Error() string { return fmt.Sprintf("Session interrompue (%s)", e.Signal) }

// RunSession supervises only the newly created session; existing ADB processes
// are never stopped. Signals are registered before starting the child.
func RunSession(cmd *exec.Cmd) error {
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	return runSession(cmd, signals)
}
