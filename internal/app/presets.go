package app

import (
	"errors"
	"strings"

	"github.com/mattn/go-shellwords"
)

func Defaults() Config {
	return Config{Presets: []Preset{
		{"Léger Wi-Fi", "800p · 2 Mbit/s · 30 FPS · sans audio", []string{"-m", "800", "-b", "2M", "--max-fps=30", "--video-codec=h264", "--no-audio"}},
		{"Très léger", "640p · 1 Mbit/s · 25 FPS · sans audio", []string{"-m", "640", "-b", "1M", "--max-fps=25", "--no-audio"}},
		{"Qualité", "1920p · 12 Mbit/s · 60 FPS", []string{"-m", "1920", "-b", "12M", "--max-fps=60"}},
	}}
}
func ParseArgs(s string) ([]string, error) {
	p := shellwords.NewParser()
	p.ParseEnv = false
	p.ParseBacktick = false
	return p.Parse(s)
}
func ValidatePreset(p Preset) error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("le nom est obligatoire")
	}
	for _, s := range append([]string{p.Name, p.Description}, p.Args...) {
		if strings.ContainsAny(s, "\x00\x1b\r\n") {
			return errors.New("caractères de contrôle interdits")
		}
	}
	return ValidateArgs(p.Args)
}
func ValidateArgs(args []string) error {
	for _, s := range args {
		if strings.ContainsAny(s, "\x00\x1b\r\n") {
			return errors.New("argument invalide")
		}
		if strings.HasPrefix(s, "--") {
			name := strings.SplitN(s, "=", 2)[0]
			if len(name) > 2 {
				for _, reserved := range []string{"--serial", "--select-usb", "--select-tcpip", "--tcpip"} {
					if strings.HasPrefix(reserved, name) {
						return errors.New("la sélection de cible est réservée à l'application")
					}
				}
			}
		} else if strings.HasPrefix(s, "-") {
			// getopt short clusters: walk known flags until an option consumes
			// the remaining characters as its value (or an unknown option).
			for _, c := range strings.TrimPrefix(s, "-") {
				if c == 's' || c == 'd' || c == 'e' {
					return errors.New("la sélection de cible est réservée à l'application")
				}
				if !strings.ContainsRune("fGhKMnNStvwx", c) {
					break
				}
			}
		}
		if s == "-s" || strings.HasPrefix(s, "-s") && !strings.HasPrefix(s, "--") || s == "-d" || s == "-e" || s == "--" || s == "--serial" || strings.HasPrefix(s, "--serial=") || s == "--select-usb" || s == "--select-tcpip" || s == "--tcpip" || strings.HasPrefix(s, "--tcpip=") {
			return errors.New("la sélection de cible est réservée à l'application")
		}
	}
	return nil
}
