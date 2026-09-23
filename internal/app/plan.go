package app

import (
	"context"
	"errors"
	"os/exec"
	"strings"
)

func (a *App) Plan(serial string, p Preset, extra []string) ([]string, error) {
	if serial == "" {
		return nil, errors.New("aucun appareil utilisable sélectionné")
	}
	if err := ValidatePreset(p); err != nil {
		return nil, err
	}
	if err := ValidateArgs(extra); err != nil {
		return nil, err
	}
	args := []string{"-s", serial}
	args = append(args, p.Args...)
	return append(args, extra...), nil
}
func (a *App) Prepare(ctx context.Context, serial string, p Preset, extra []string) (*exec.Cmd, error) {
	args, err := a.Plan(serial, p, extra)
	if err != nil {
		return nil, err
	}
	executable, err := exec.LookPath("scrcpy")
	if err != nil {
		return nil, errors.New("scrcpy introuvable : installez-le et vérifiez PATH")
	}
	if err := a.checkCapabilities(ctx, executable, p); err != nil {
		return nil, err
	}
	ds, err := a.Discover(ctx)
	if err != nil {
		return nil, err
	}
	for _, d := range ds {
		if d.Serial == serial && d.State == "device" {
			return exec.Command(executable, args...), nil
		}
	}
	return nil, errors.New("l'appareil n'est plus disponible ; rafraîchissez avec r")
}
func Quote(s string) string {
	if s != "" && !strings.ContainsAny(s, " \t\n'\"\\$`;|&()<>{}*?!") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
func Command(args []string) string {
	parts := []string{"scrcpy"}
	for _, a := range args {
		parts = append(parts, Quote(a))
	}
	return strings.Join(parts, " ")
}
