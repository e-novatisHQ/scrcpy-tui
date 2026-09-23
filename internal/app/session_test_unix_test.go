//go:build !windows

package app

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"testing"
)

func ensureSessionTestConsole(t *testing.T) {}

func watchSessionProcess(t *testing.T, pid int) func() bool {
	t.Helper()
	return func() bool {
		if syscall.Kill(pid, 0) == syscall.ESRCH {
			return false
		}
		// Linux containers can retain a dead orphan as a zombie until PID 1 reaps it.
		if data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid)); err == nil {
			end := strings.LastIndex(string(data), ")")
			if end >= 0 && strings.HasPrefix(string(data[end+1:]), " Z") {
				return false
			}
		}
		return true
	}
}
