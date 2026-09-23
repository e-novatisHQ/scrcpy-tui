package app

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Runner interface {
	Output(context.Context, string, ...string) ([]byte, error)
}
type System struct{}

const maxToolOutput = 1024 * 1024

type toolOutput struct {
	data      []byte
	truncated bool
}

func (b *toolOutput) Write(p []byte) (int, error) {
	n := min(len(p), maxToolOutput-len(b.data))
	b.data = append(b.data, p[:n]...)
	b.truncated = b.truncated || n != len(p)
	return len(p), nil
}

func (System) Output(ctx context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	// A wrapper's descendant retaining a pipe must not defeat the probe timeout.
	cmd.WaitDelay = 2 * time.Second
	var output toolOutput
	// os/exec uses one copy goroutine when these comparable writers are identical.
	cmd.Stdout, cmd.Stderr = &output, &output
	err := cmd.Run()
	if output.truncated {
		err = errors.Join(err, errors.New("sortie de commande trop volumineuse (limite 1 Mio)"))
	}
	return output.data, err
}

func ParseDevices(b []byte) []Device {
	var ds []Device
	for _, line := range strings.Split(string(b), "\n") {
		f := strings.Fields(line)
		if len(f) < 2 || strings.HasPrefix(line, "List of devices") || strings.HasPrefix(line, "*") {
			continue
		}
		d := Device{Serial: f[0], State: f[1]}
		for _, v := range f[2:] {
			if strings.HasPrefix(v, "model:") {
				d.Model = strings.ReplaceAll(strings.TrimPrefix(v, "model:"), "_", " ")
			}
		}
		ds = append(ds, d)
	}
	return ds
}
func (a *App) Discover(ctx context.Context) ([]Device, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	b, err := a.Runner.Output(ctx, "adb", "devices", "-l")
	if err != nil {
		return nil, fmt.Errorf("échec ADB : %w (%s)", err, strings.TrimSpace(string(b)))
	}
	return ParseDevices(b), nil
}
