package app

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type Interrupted struct{ Signal os.Signal }

func (e *Interrupted) Error() string { return fmt.Sprintf("Session interrompue (%s)", e.Signal) }

// RunSession owns only the newly created process group; existing ADB processes
// are never stopped. Signals are registered before starting the child.
func RunSession(cmd *exec.Cmd) error {
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(signals)
	return runSession(cmd, signals)
}
func runSession(cmd *exec.Cmd, signals <-chan os.Signal) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	// A descendant retaining an output pipe must not hold Wait forever.
	cmd.WaitDelay = 2 * time.Second
	if err := cmd.Start(); err != nil {
		return err
	}
	group := cmd.Process.Pid
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		stopGroup(group)
		return err
	case sig := <-signals:
		syscall.Kill(-group, syscall.SIGTERM)
		until := time.Now().Add(2 * time.Second)
		deadline := time.NewTimer(time.Until(until))
		defer deadline.Stop()
		select {
		case <-done:
			stopGroupUntil(group, until)
		case <-deadline.C:
			syscall.Kill(-group, syscall.SIGKILL)
			<-done
		}
		return &Interrupted{Signal: sig}
	}
}
func stopGroup(group int) { stopGroupUntil(group, time.Now().Add(2*time.Second)) }
func stopGroupUntil(group int, until time.Time) {
	// Wait for graceful termination of any descendants after the parent exits.
	if syscall.Kill(-group, 0) == syscall.ESRCH {
		return
	}
	syscall.Kill(-group, syscall.SIGTERM)
	for time.Now().Before(until) {
		if syscall.Kill(-group, 0) == syscall.ESRCH {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	syscall.Kill(-group, syscall.SIGKILL)
}
