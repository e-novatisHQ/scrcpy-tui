package app

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Runner interface {
	Output(context.Context, string, ...string) ([]byte, error)
}
type System struct{}

func (System) Output(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
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
