package app

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
)

// Match declarations, not references in the more deeply indented prose. Both
// long options and their short aliases are needed by the bundled presets.
var optionDeclaration = regexp.MustCompile(`(?m)^[ \t]{1,4}(?:(-[A-Za-z0-9?]),[ \t]+)?(--[a-z][a-z0-9-]*|-[A-Za-z0-9?])(?:[ =,\t\r]|$)`)

func parseCapabilities(help []byte) map[string]bool {
	options := make(map[string]bool)
	for _, match := range optionDeclaration.FindAllStringSubmatch(string(help), -1) {
		options[match[2]] = true
		if match[1] != "" {
			options[match[1]] = true
		}
	}
	return options
}

// Only unchanged bundled argv have a known grammar here. User-authored presets
// and free --args retain their existing contract: scrcpy validates them. Names
// are not identifiers, so renaming a bundled preset does not bypass the check.
func bundledPreset(p Preset) bool {
	for _, builtin := range Defaults().Presets {
		if slices.Equal(p.Args, builtin.Args) {
			return true
		}
	}
	return false
}

func (a *App) checkCapabilities(ctx context.Context, executable string, p Preset) error {
	if !bundledPreset(p) {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	help, err := a.Runner.Output(ctx, executable, "--help")
	if err != nil {
		return fmt.Errorf("capacités scrcpy indisponibles : %w ; vérifiez scrcpy --help", err)
	}
	if len(help) > maxToolOutput {
		return fmt.Errorf("aide scrcpy trop volumineuse ; vérifiez scrcpy --help")
	}
	available := parseCapabilities(help)
	if len(available) == 0 {
		return fmt.Errorf("format d'aide scrcpy non reconnu ; vérifiez scrcpy --help")
	}
	required := map[string]bool{"-s": true}
	for _, arg := range p.Args {
		if strings.HasPrefix(arg, "-") {
			name, _, _ := strings.Cut(arg, "=")
			required[name] = true
		}
	}
	var missing []string
	for option := range required {
		if !available[option] {
			missing = append(missing, option)
		}
	}
	if len(missing) != 0 {
		sort.Strings(missing)
		return fmt.Errorf("options scrcpy indisponibles pour le preset %q : %s ; mettez à jour scrcpy ou adaptez le preset (scrcpy --help)", p.Name, strings.Join(missing, ", "))
	}
	return nil
}
