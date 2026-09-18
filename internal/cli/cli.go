// Package cli provides the testable command-line entrypoint.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"

	"github.com/e-novatisHQ/scrcpy-tui/internal/app"
	"github.com/e-novatisHQ/scrcpy-tui/internal/tui"
)

type Streams struct {
	In       *os.File
	Out, Err io.Writer
}

func Run(argsInput []string, version string, streams Streams) int {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir = "."
	}
	fs := flag.NewFlagSet("scrcpy-tui", flag.ContinueOnError)
	fs.SetOutput(streams.Err)
	config := fs.String("config", filepath.Join(dir, "scrcpy-tui", "presets.json"), "Fichier de configuration")
	device := fs.String("device", "", "Identifiant ADB")
	preset := fs.String("preset", "Léger Wi-Fi", "Nom du preset")
	extra := fs.String("args", "", "Arguments scrcpy supplémentaires (guillemets acceptés)")
	yes := fs.Bool("yes", false, "Autoriser le lancement CLI")
	ver := fs.Bool("version", false, "Afficher la version")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: scrcpy-tui [options] [tui|devices|presets|preview|launch]\nOptions avant ou après la sous-commande ; sans commande : TUI.")
		fs.PrintDefaults()
	}
	if err := fs.Parse(normalizeArgs(argsInput)); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 3
	}
	if *ver {
		fmt.Fprintln(streams.Out, version)
		return 0
	}
	command := "tui"
	if fs.NArg() > 0 {
		command = fs.Arg(0)
	}
	if fs.NArg() > 1 {
		fmt.Fprintln(streams.Err, "Arguments inattendus :", fs.Args()[1:])
		fs.Usage()
		return 3
	}
	a, warning := app.New(*config)
	if warning != "" {
		fmt.Fprintln(streams.Err, warning)
	}
	args, err := app.ParseArgs(*extra)
	if err != nil {
		fmt.Fprintln(streams.Err, err)
		return 3
	}
	switch command {
	case "tui":
		if !term.IsTerminal(int(streams.In.Fd())) {
			fmt.Fprintln(streams.Err, "La TUI exige un terminal ; utilisez devices, presets, preview ou launch.")
			return 3
		}
		m := tui.New(a, warning)
		m.Extra = args
		_, err = tea.NewProgram(m, tea.WithAltScreen(), tea.WithInput(streams.In), tea.WithOutput(streams.Out)).Run()
	case "devices":
		var ds []app.Device
		ds, err = a.Discover(context.Background())
		for _, d := range ds {
			fmt.Fprintf(streams.Out, "%s\t%s\t%s\n", d.Serial, d.State, d.Model)
		}
	case "presets":
		for _, p := range a.Config.Presets {
			fmt.Fprintf(streams.Out, "%s\t%s\n", p.Name, p.Description)
		}
	case "preview", "launch":
		var p app.Preset
		found := false
		for _, v := range a.Config.Presets {
			if v.Name == *preset {
				p = v
				found = true
			}
		}
		if !found {
			fmt.Fprintln(streams.Err, "Preset inconnu")
			return 3
		}
		var plan []string
		plan, err = a.Plan(*device, p, args)
		if err != nil {
			fmt.Fprintln(streams.Err, err)
			return 3
		}
		fmt.Fprintln(streams.Out, app.Command(plan))
		if command == "preview" {
			return 0
		}
		if !*yes {
			fmt.Fprintln(streams.Err, "Lancement CLI : --yes requis")
			return 4
		}
		cmd, e := a.Prepare(context.Background(), *device, p, args)
		if e != nil {
			err = e
			break
		}
		a.Config.LastDevice = *device
		a.Config.LastPreset = p.Name
		if e = a.Save(); e != nil {
			fmt.Fprintln(streams.Err, "Configuration :", e)
		}
		cmd.Stdin = streams.In
		cmd.Stdout = streams.Out
		cmd.Stderr = streams.Err
		err = app.RunSession(cmd)
		var interrupted *app.Interrupted
		if errors.As(err, &interrupted) {
			return 130
		}

	default:
		fs.Usage()
		return 3
	}
	if err != nil {
		fmt.Fprintln(streams.Err, err)
		return 1
	}
	return 0
}

// Preserve option values verbatim, including values which resemble commands.
func normalizeArgs(args []string) []string {
	var flags, positional []string
	terminated := false
	for i := 0; i < len(args); i++ {
		s := args[i]
		if s == "--" {
			terminated = true
			positional = append(positional, args[i+1:]...)
			break
		}
		if strings.HasPrefix(s, "-") {
			flags = append(flags, s)
			name := strings.TrimLeft(s, "-")
			if !strings.Contains(name, "=") && (name == "config" || name == "device" || name == "preset" || name == "args") && i+1 < len(args) {
				i++
				flags = append(flags, args[i])
			}
		} else {
			positional = append(positional, s)
		}
	}
	if terminated {
		flags = append(flags, "--")
	}
	return append(flags, positional...)
}
