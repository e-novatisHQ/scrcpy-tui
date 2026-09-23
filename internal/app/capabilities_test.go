package app

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/e-novatisHQ/scrcpy-tui/internal/testutil"
)

const completeHelp = `scrcpy test-build
Options:
    -s, --serial=serial
    -m, --max-size=value
    -b, --video-bit-rate=value
    --max-fps=value
    --video-codec=name
    --no-audio
`

type capabilityRunner struct {
	help  []byte
	err   error
	wait  bool
	calls []string
}

func (r *capabilityRunner) Output(ctx context.Context, name string, args ...string) ([]byte, error) {
	r.calls = append(r.calls, strings.Join(args, " "))
	if slices.Equal(args, []string{"--help"}) {
		if r.wait {
			<-ctx.Done()
			return nil, ctx.Err()
		}
		return r.help, r.err
	}
	return []byte("USB device model:Fake_TV\n"), nil
}

func TestCapabilityDeclarations(t *testing.T) {
	help := completeHelp + "        This description mentions --unknown.\n        --not-a-declaration is prose.\n"
	options := parseCapabilities([]byte(help))
	for _, name := range []string{"-s", "-m", "-b", "--serial", "--max-size", "--video-bit-rate", "--max-fps", "--video-codec", "--no-audio"} {
		if !options[name] {
			t.Errorf("missing declaration %s", name)
		}
	}
	if options["--unknown"] || options["--not-a-declaration"] {
		t.Fatal("prose accepted as an option")
	}
	if !parseCapabilities([]byte("  -m pixels\r\n  --no-audio\r\n"))["--no-audio"] {
		t.Fatal("CRLF help not parsed")
	}
}

func TestCapabilityChecksSelectedBuiltin(t *testing.T) {
	for _, p := range Defaults().Presets {
		r := &capabilityRunner{help: []byte(completeHelp)}
		a := &App{Runner: r}
		p.Name = "Renamed builtin"
		if err := a.checkCapabilities(context.Background(), "scrcpy", p); err != nil {
			t.Fatal(err)
		}
	}
	missing := strings.ReplaceAll(completeHelp, "    --no-audio\n", "")
	r := &capabilityRunner{help: []byte(missing)}
	a := &App{Runner: r}
	err := a.checkCapabilities(context.Background(), "scrcpy", Defaults().Presets[0])
	if err == nil || !strings.Contains(err.Error(), "--no-audio") || !strings.Contains(err.Error(), "adaptez") {
		t.Fatal(err)
	}
	if err := a.checkCapabilities(context.Background(), "scrcpy", Defaults().Presets[2]); err != nil {
		t.Fatal("unneeded option blocked selected preset", err)
	}
	r.calls = nil
	if err := a.checkCapabilities(context.Background(), "scrcpy", Preset{Name: "Custom", Args: []string{"--future-option"}}); err != nil || len(r.calls) != 0 {
		t.Fatal("custom argv contract changed", err, r.calls)
	}
}

func TestCapabilitiesFailClosed(t *testing.T) {
	for _, r := range []*capabilityRunner{
		{help: []byte("description mentioning --no-audio, but no option declarations")},
		{err: errors.New("tool failed")},
		{help: []byte(strings.Repeat("x", 1024*1024+1))},
		{wait: true},
	} {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
		err := (&App{Runner: r}).checkCapabilities(ctx, "scrcpy", Defaults().Presets[0])
		cancel()
		if err == nil {
			t.Fatal("unverifiable help accepted")
		}
	}
}

func TestPrepareChecksCapabilitiesBeforeDeviceContact(t *testing.T) {
	dir := t.TempDir()
	executable := testutil.Copy(t, testutil.Build(t, "helper"), dir, "scrcpy")
	t.Setenv("PATH", dir)
	r := &capabilityRunner{help: []byte(strings.ReplaceAll(completeHelp, "    --no-audio\n", ""))}
	a := &App{Runner: r}
	if _, err := a.Prepare(context.Background(), "USB", Defaults().Presets[0], nil); err == nil {
		t.Fatal("incompatible builtin prepared")
	}
	if !slices.Equal(r.calls, []string{"--help"}) {
		t.Fatal("device contacted on incompatible preset", r.calls)
	}
	r.help = []byte(completeHelp)
	r.calls = nil
	cmd, err := a.Prepare(context.Background(), "USB", Defaults().Presets[0], []string{"--unknown=user-choice"})
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Path != executable || !slices.Equal(r.calls, []string{"--help", "devices -l"}) {
		t.Fatal(cmd.Path, r.calls)
	}
	if _, err := os.Stat(filepath.Join(dir, "profile.json")); !os.IsNotExist(err) {
		t.Fatal("prepare wrote configuration")
	}
}

func TestSystemBoundsToolOutput(t *testing.T) {
	helper := testutil.Build(t, "probe helper é")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	output, err := (System{}).Output(ctx, helper, "flood")
	if err == nil || !strings.Contains(err.Error(), "volumineuse") || len(output) != maxToolOutput {
		t.Fatal(len(output), err)
	}
}
