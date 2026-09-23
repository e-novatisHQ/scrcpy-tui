//go:build !windows

package app

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestSessionSignalsAreRelayed(t *testing.T) {
	for _, sig := range []os.Signal{os.Interrupt, syscall.SIGTERM} {
		t.Run(sig.String(), func(t *testing.T) {
			ready := filepath.Join(t.TempDir(), "ready")
			cmd := exec.Command("sh", "-c", `trap 'exit 0' TERM; echo ready > "$1"; while :; do sleep 1; done`, "fake", ready)
			signals := make(chan os.Signal, 1)
			result := make(chan error, 1)
			go func() { result <- runSession(cmd, signals) }()
			deadline := time.Now().Add(3 * time.Second)
			for {
				if _, err := os.Stat(ready); err == nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("not started")
				}
				time.Sleep(10 * time.Millisecond)
			}
			signals <- sig
			select {
			case err := <-result:
				var interrupted *Interrupted
				if !errors.As(err, &interrupted) || interrupted.Signal != sig {
					t.Fatal(err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("signal did not stop session")
			}
		})
	}
}
func TestSessionOrdinaryExit(t *testing.T) {
	if err := runSession(exec.Command("sh", "-c", "exit 0"), make(chan os.Signal)); err != nil {
		t.Fatal(err)
	}
	if err := runSession(exec.Command("sh", "-c", "exit 7"), make(chan os.Signal)); err == nil {
		t.Fatal("exit failure ignored")
	}
}

func TestSessionForcesUnresponsiveOwnedProcess(t *testing.T) {
	ready := filepath.Join(t.TempDir(), "ready")
	cmd := exec.Command("sh", "-c", `trap '' TERM; echo ready > "$1"; while :; do sleep 1; done`, "fake", ready)
	signals := make(chan os.Signal, 1)
	result := make(chan error, 1)
	go func() { result <- runSession(cmd, signals) }()
	limit := time.Now().Add(3 * time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(limit) {
			t.Fatal("not started")
		}
		time.Sleep(10 * time.Millisecond)
	}
	signals <- syscall.SIGTERM
	select {
	case err := <-result:
		var interrupted *Interrupted
		if !errors.As(err, &interrupted) {
			t.Fatal(err)
		}
	case <-time.After(4 * time.Second):
		t.Fatal("forced stop not bounded")
	}
}

func TestExitedParentWithInheritedOutputDoesNotHang(t *testing.T) {
	cmd := exec.Command("sh", "-c", "sleep 20 & exit 0")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	result := make(chan error, 1)
	go func() { result <- runSession(cmd, make(chan os.Signal)) }()
	select {
	case err := <-result:
		if !errors.Is(err, exec.ErrWaitDelay) {
			t.Fatal("inherited output not diagnosed", err)
		}
	case <-time.After(6 * time.Second):
		t.Fatal("descendant output keeps session hanging")
	}
}
