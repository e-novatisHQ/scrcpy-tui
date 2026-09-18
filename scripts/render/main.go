// Generate plain-text fixtures for visual review without ADB or profile writes.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/x/ansi"

	"github.com/e-novatisHQ/scrcpy-tui/internal/app"
	"github.com/e-novatisHQ/scrcpy-tui/internal/tui"
)

func main() {
	a := &app.App{Config: app.Defaults()}
	m := tui.New(a, "")
	m.D = 0
	m.Busy = false
	for i := 0; i < 42; i++ {
		m.Devices = append(m.Devices, app.Device{Serial: fmt.Sprintf("192.0.2.10:%d", 6201+i), State: "device", Model: "Dongle G 4K"})
	}
	for _, s := range []struct {
		name string
		w, h int
	}{{"narrow", 40, 14}, {"standard", 100, 24}, {"large", 140, 48}} {
		m.W = s.w
		m.H = s.h
		path := filepath.Join("docs", "renderings", s.name+".txt")
		if err := os.WriteFile(path, []byte(plain(m.View())+"\n"), 0644); err != nil {
			panic(err)
		}
	}
}

func plain(view string) string {
	rows := strings.Split(ansi.Strip(view), "\n")
	for i, s := range rows {
		rows[i] = strings.TrimRight(s, " \t")
	}
	return strings.Join(rows, "\n")
}
