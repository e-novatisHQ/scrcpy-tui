package app

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/e-novatisHQ/scrcpy-tui/internal/testutil"
)

type fake struct {
	b   []byte
	err error
}

func (f fake) Output(context.Context, string, ...string) ([]byte, error) { return f.b, f.err }
func TestDevicesAndRevalidation(t *testing.T) {
	a := &App{Runner: fake{b: []byte("List of devices attached\nUSB device model:Pixel_8\nBAD unauthorized\nNET offline\n")}}
	ds, err := a.Discover(context.Background())
	if err != nil || len(ds) != 3 || ds[0].Model != "Pixel 8" {
		t.Fatal(ds, err)
	}
	a.Runner = fake{b: []byte("USB offline\n")}
	t.Setenv("PATH", t.TempDir())
	if _, err = a.Prepare(context.Background(), "USB", Defaults().Presets[0], nil); err == nil {
		t.Fatal("missing dependency accepted")
	}
}
func TestArguments(t *testing.T) {
	args, err := ParseArgs(`--window-title 'Une TV' --unknown='a b' '$HOME'`)
	if err != nil || !reflect.DeepEqual(args, []string{"--window-title", "Une TV", "--unknown=a b", "$HOME"}) {
		t.Fatal(args, err)
	}
	a := &App{}
	plan, err := a.Plan("TV", Preset{Name: "Test", Args: args}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if plan[0] != "-s" || plan[1] != "TV" {
		t.Fatal(plan)
	}
	for _, s := range []string{"-sOTHER", "--serial=x", "--select-usb", "--tcpip=host"} {
		if ValidateArgs([]string{s}) == nil {
			t.Fatal(s)
		}
	}
	if _, err = ParseArgs(`'broken`); err == nil {
		t.Fatal("quote accepted")
	}
}
func TestPersistenceAndRecovery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config", "presets.json")
	a, w := New(path)
	if w != "" || len(a.Config.Presets) != 3 {
		t.Fatal(w)
	}
	a.Config.Presets = append(a.Config.Presets, Preset{Name: "Perso", Args: []string{"--unknown"}})
	a.Config.LastDevice = "USB"
	if err := a.Save(); err != nil {
		t.Fatal(err)
	}
	b, w := New(path)
	if w != "" || len(b.Config.Presets) != 4 || b.Config.LastDevice != "USB" {
		t.Fatal(w, b)
	}
	st, _ := os.Stat(path)
	if st.Mode().Perm() != 0600 {
		t.Fatal(st.Mode())
	}
	os.WriteFile(path, []byte("{"), 0600)
	_, w = New(path)
	if w == "" {
		t.Fatal("no warning")
	}
	raw, _ := os.ReadFile(path)
	if string(raw) != "{" {
		t.Fatal("overwritten")
	}
}

func TestPrepareRejectsDisconnectedDevice(t *testing.T) {
	dir := t.TempDir()
	testutil.Copy(t, testutil.Build(t, "helper"), dir, "scrcpy")
	t.Setenv("PATH", dir)
	a := &App{Runner: fake{b: []byte("USB device model:TV\n")}}
	p := Defaults().Presets[0]
	cmd, err := a.Prepare(context.Background(), "USB", p, nil)
	if err != nil || cmd.Args[1] != "-s" || cmd.Args[2] != "USB" {
		t.Fatal(cmd, err)
	}
	for _, state := range []string{"offline", "unauthorized"} {
		a.Runner = fake{b: []byte("USB " + state + "\n")}
		if _, err = a.Prepare(context.Background(), "USB", p, nil); err == nil {
			t.Fatal("disabled device accepted", state)
		}
	}
	a.Runner = fake{}
	if _, err = a.Prepare(context.Background(), "USB", p, nil); err == nil {
		t.Fatal("missing device accepted")
	}
}

func TestRecoveryPreservesOriginalOnSave(t *testing.T) {
	for _, original := range []string{"{broken", `{"presets":[{"name":"Keep","args":[]},{"name":"","args":[]}],"custom":"preserve"}`} {
		t.Run(original, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "presets.json")
			os.WriteFile(path, []byte(original), 0600)
			a, warning := New(path)
			if warning == "" {
				t.Fatal("missing recovery warning")
			}
			a.Config.LastDevice = "USB"
			if err := a.Save(); err != nil {
				t.Fatal(err)
			}
			backups, _ := filepath.Glob(path + ".recovery-*")
			if len(backups) != 1 {
				t.Fatal("original not backed up", backups)
			}
			data, _ := os.ReadFile(backups[0])
			if string(data) != original {
				t.Fatal("original changed")
			}
			st, _ := os.Stat(backups[0])
			if st.Mode().Perm() != 0600 {
				t.Fatal("backup permission")
			}
			a.Save()
			backups, _ = filepath.Glob(path + ".recovery-*")
			if len(backups) != 1 {
				t.Fatal("repeated backup")
			}
		})
	}
}

func TestUnreadableConfigNeverOverwritten(t *testing.T) {
	path := t.TempDir()
	a, warning := New(path)
	if warning == "" {
		t.Fatal("directory accepted as config")
	}
	if err := a.Save(); err == nil {
		t.Fatal("unreadable config overwritten")
	}
	st, err := os.Stat(path)
	if err != nil || !st.IsDir() {
		t.Fatal("directory changed")
	}
}

func TestSelectorsCannotHideInShortClustersOrAbbreviations(t *testing.T) {
	for _, arg := range []string{"-fsOTHER", "-fnsOTHER", "-fd", "-ne", "--ser=OTHER", "--select-u", "--tcp=OTHER"} {
		if err := ValidateArgs([]string{arg}); err == nil {
			t.Fatal("hidden selector accepted", arg)
		}
	}
	for _, arg := range []string{"-m800", "-b2M", "-Vdebug", "-rside.mp4", "-fnm800", "--custom-unknown=value"} {
		if err := ValidateArgs([]string{arg}); err != nil {
			t.Fatal("valid argument rejected", arg, err)
		}
	}
}
