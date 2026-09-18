package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func invoke(t *testing.T, args ...string) (int, string, string) {
	t.Helper()
	in, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	var out, stderr bytes.Buffer
	code := Run(args, "test-version", Streams{In: in, Out: &out, Err: &stderr})
	return code, out.String(), stderr.String()
}
func TestEntrypointReadOnlyCommands(t *testing.T) {
	config := filepath.Join(t.TempDir(), "presets.json")
	for _, tt := range []struct {
		args     []string
		code     int
		out, err string
	}{{[]string{"--version"}, 0, "test-version", ""}, {[]string{"--help"}, 0, "", "Usage:"}, {[]string{"presets"}, 0, "Léger Wi-Fi", ""}, {[]string{"preview", "--device", "USB"}, 0, "scrcpy -s USB", ""}, {[]string{"preview", "--device", "USB", "--preset", "absent"}, 3, "", "Preset inconnu"}, {[]string{"launch", "--device", "USB"}, 4, "scrcpy -s USB", "--yes requis"}, {[]string{"preview", "--device", "USB", "--args=-fsOTHER"}, 3, "", "réservée"}, {[]string{"tui"}, 3, "", "terminal"}, {[]string{"absent"}, 3, "", "Usage:"}, {[]string{"preview", "--", "--yes"}, 3, "", "Arguments inattendus"}} {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			args := append([]string{"--config", config}, tt.args...)
			code, out, err := invoke(t, args...)
			if code != tt.code || !strings.Contains(out, tt.out) || !strings.Contains(err, tt.err) {
				t.Fatal(code, out, err)
			}
			if _, err := os.Stat(config); !os.IsNotExist(err) {
				t.Fatal("read-only command wrote a config")
			}
		})
	}
}
func TestEntrypointDiscoveryAndLaunch(t *testing.T) {
	dir := t.TempDir()
	for name, body := range map[string]string{"adb": "#!/bin/sh\nprintf 'List of devices attached\\nUSB device model:Fake_TV\\nBAD unauthorized\\n'\n", "scrcpy": "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$CAPTURE_PATH\"\n"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0700); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("PATH", dir)
	capture := filepath.Join(dir, "args")
	t.Setenv("CAPTURE_PATH", capture)
	config := filepath.Join(dir, "profile.json")
	code, out, err := invoke(t, "devices", "--config", config)
	if code != 0 || !strings.Contains(out, "Fake TV") || err != "" {
		t.Fatal(code, out, err)
	}
	code, out, err = invoke(t, "launch", "--config", config, "--device", "BAD", "--yes")
	if code != 1 {
		t.Fatal(code, out, err)
	}
	if _, err := os.Stat(capture); !os.IsNotExist(err) {
		t.Fatal("disabled device launched")
	}
	code, out, err = invoke(t, "launch", "--config", config, "--device", "USB", "--preset", "Très léger", "--yes")
	if code != 0 {
		t.Fatal(code, out, err)
	}
	args, e := os.ReadFile(capture)
	if e != nil || !strings.HasPrefix(string(args), "-s\nUSB\n") || !strings.Contains(string(args), "--max-fps=25") {
		t.Fatal(string(args), e)
	}
	data, e := os.ReadFile(config)
	if e != nil || !strings.Contains(string(data), `"last_device": "USB"`) {
		t.Fatal(string(data), e)
	}
}
func TestEntrypointMissingTools(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	code, _, err := invoke(t, "devices", "--config", filepath.Join(t.TempDir(), "profile.json"))
	if code != 1 || !strings.Contains(err, "ADB") {
		t.Fatal(code, err)
	}
}
