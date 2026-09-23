package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type App struct {
	recoveryBytes []byte
	loadErr       error

	Config Config
	Path   string
	Runner Runner
}

func New(path string) (*App, string) {
	a := &App{Config: Defaults(), Path: path, Runner: System{}}
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return a, ""
	}
	if err != nil {
		a.loadErr = err
		return a, "Configuration inaccessible : " + err.Error()
	}
	var c Config
	if err = json.Unmarshal(b, &c); err != nil {
		a.recoveryBytes = b
		return a, "Configuration invalide ; presets par défaut chargés. Copie de secours avant enregistrement."
	}
	valid := []Preset{}
	seen := map[string]bool{}
	bad := false
	for _, p := range c.Presets {
		if ValidatePreset(p) != nil || seen[p.Name] {
			bad = true
			continue
		}
		seen[p.Name] = true
		valid = append(valid, p)
	}
	c.Presets = valid
	a.Config = c
	if bad {
		a.recoveryBytes = b
		return a, "Certains presets invalides ont été ignorés. Copie de secours avant enregistrement."
	}
	return a, ""
}
func (a *App) Save() error {
	if a.loadErr != nil {
		return fmt.Errorf("configuration illisible, fichier préservé : %w", a.loadErr)
	}
	if err := os.MkdirAll(filepath.Dir(a.Path), 0700); err != nil {
		return err
	}
	if a.recoveryBytes != nil {
		f, err := os.CreateTemp(filepath.Dir(a.Path), filepath.Base(a.Path)+".recovery-*")
		if err != nil {
			return err
		}
		if err = restrictConfigFile(f); err == nil {
			_, err = f.Write(a.recoveryBytes)
		}
		if err == nil {
			err = f.Sync()
		}
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
		if err != nil {
			os.Remove(f.Name())
			return fmt.Errorf("copie de secours impossible, fichier préservé : %w", err)
		}
		a.recoveryBytes = nil
	}
	b, err := json.MarshalIndent(a.Config, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(a.Path), ".presets-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err = restrictConfigFile(f); err == nil {
		_, err = f.Write(b)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(f.Name(), a.Path)
}
