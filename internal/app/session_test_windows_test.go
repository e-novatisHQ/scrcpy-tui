package app

import (
	"syscall"
	"testing"

	"golang.org/x/sys/windows"
)

func ensureSessionTestConsole(t *testing.T) {
	t.Helper()
	// Hosted runners may start tests without a console. Allocate one so that the
	// native graceful-control-event contract is actually exercised.
	r, _, err := windows.NewLazySystemDLL("kernel32.dll").NewProc("AllocConsole").Call()
	if r == 0 && err != syscall.ERROR_ACCESS_DENIED {
		t.Fatalf("allocate test console: %v", err)
	}
}

func watchSessionProcess(t *testing.T, pid int) func() bool {
	t.Helper()
	h, err := windows.OpenProcess(windows.SYNCHRONIZE, false, uint32(pid))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { windows.CloseHandle(h) })
	return func() bool {
		status, err := windows.WaitForSingleObject(h, 0)
		if err != nil {
			t.Fatal(err)
		}
		return status == uint32(windows.WAIT_TIMEOUT)
	}
}
