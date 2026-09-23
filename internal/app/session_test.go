package app

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/e-novatisHQ/scrcpy-tui/internal/testutil"
)

func waitSessionReady(t *testing.T, path string) int {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		b, err := os.ReadFile(path)
		if err == nil {
			pid, err := strconv.Atoi(strings.TrimSpace(string(b)))
			if err == nil && pid > 0 {
				return pid
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("helper did not become ready")
	return 0
}

func TestNativeSession(t *testing.T) {
	ensureSessionTestConsole(t)
	helper := testutil.Build(t, "session-helper")
	t.Run("exit", func(t *testing.T) {
		if err := runSession(exec.Command(helper, "exit"), make(chan os.Signal)); err != nil {
			t.Fatal(err)
		}
		err := runSession(exec.Command(helper, "fail"), make(chan os.Signal))
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 7 {
			t.Fatalf("exit status: %v", err)
		}
		if err := runSession(exec.Command(filepath.Join(t.TempDir(), "missing")), make(chan os.Signal)); err == nil {
			t.Fatal("missing executable accepted")
		}
	})
	for _, mode := range []string{"wait", "ignore", "descendant", "orphan"} {
		t.Run(mode, func(t *testing.T) {
			ready := filepath.Join(t.TempDir(), "ready")
			// A preexisting, unrelated process must survive every cleanup path.
			outsideReady := filepath.Join(t.TempDir(), "outside")
			outside := exec.Command(helper, "ignore", outsideReady)
			if err := outside.Start(); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = outside.Process.Kill(); _ = outside.Wait() })
			outsideAlive := watchSessionProcess(t, waitSessionReady(t, outsideReady))
			cmd := exec.Command(helper, mode, ready)
			cmd.Stdout, cmd.Stderr = io.Discard, io.Discard
			signals := make(chan os.Signal, 1)
			result := make(chan error, 1)
			go func() { result <- runSession(cmd, signals) }()
			t.Cleanup(func() {
				select {
				case signals <- os.Interrupt:
				default:
				}
			})
			pid := waitSessionReady(t, ready)
			var alive func() bool
			if mode != "orphan" {
				alive = watchSessionProcess(t, pid)
			}
			var childAlive func() bool
			if mode == "descendant" || mode == "orphan" {
				childAlive = watchSessionProcess(t, waitSessionReady(t, ready+".child"))
			}
			if mode != "orphan" {
				signals <- os.Interrupt
			}
			select {
			case err := <-result:
				if mode == "orphan" {
					if !errors.Is(err, exec.ErrWaitDelay) {
						t.Fatalf("inherited pipes: %v", err)
					}
				} else {
					var interrupted *Interrupted
					if !errors.As(err, &interrupted) || interrupted.Signal != os.Interrupt {
						t.Fatalf("interrupt: %v", err)
					}
				}
			case <-time.After(8 * time.Second):
				t.Fatal("cleanup did not finish")
			}
			deadline := time.Now().Add(3 * time.Second)
			for (alive != nil && alive()) || (childAlive != nil && childAlive()) {
				if time.Now().After(deadline) {
					t.Fatal("owned process still alive")
				}
				time.Sleep(10 * time.Millisecond)
			}
			if !outsideAlive() {
				t.Fatal("unrelated process was stopped")
			}
			if mode == "wait" {
				if _, err := os.Stat(ready + ".stopped"); err != nil {
					t.Fatalf("graceful shutdown not observed: %v", err)
				}
			}
		})
	}
}
