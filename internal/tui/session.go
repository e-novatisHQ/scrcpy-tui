package tui

import (
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/e-novatisHQ/scrcpy-tui/internal/app"
)

// tailBuffer retains at most 32 KiB, including output from both process streams.
type tailBuffer struct {
	mu sync.Mutex
	b  []byte
}

func (t *tailBuffer) Write(p []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.b = append(t.b, p...)
	if len(t.b) > 32768 {
		t.b = append([]byte{}, t.b[len(t.b)-32768:]...)
	}
	return len(p), nil
}
func (t *tailBuffer) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return strings.ToValidUTF8(string(t.b), "�")
}

type sessionCommand struct {
	cmd  *exec.Cmd
	tail tailBuffer
}

func (s *sessionCommand) SetStdin(r io.Reader)  { s.cmd.Stdin = r }
func (s *sessionCommand) SetStdout(w io.Writer) { s.cmd.Stdout = io.MultiWriter(w, &s.tail) }
func (s *sessionCommand) SetStderr(w io.Writer) { s.cmd.Stderr = io.MultiWriter(w, &s.tail) }
func (s *sessionCommand) Run() error            { return app.RunSession(s.cmd) }
