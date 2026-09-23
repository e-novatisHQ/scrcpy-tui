// A disposable native executable for integration tests; never contacts devices.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	switch strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe") {
	case "adb":
		fmt.Print("List of devices attached\nUSB device model:Fake_TV\nBAD unauthorized\n")
		return
	case "scrcpy":
		if err := os.WriteFile(os.Getenv("CAPTURE_PATH"), []byte(strings.Join(os.Args[1:], "\n")+"\n"), 0600); err != nil {
			panic(err)
		}
		return
	}
	if len(os.Args) < 2 {
		os.Exit(2)
	}
	switch os.Args[1] {
	case "exit":
		return
	case "fail":
		os.Exit(7)
	case "wait", "ignore", "descendant", "orphan":
		signals := make(chan os.Signal, 2)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		if os.Args[1] == "descendant" || os.Args[1] == "orphan" {
			child := exec.Command(os.Args[0], "ignore", os.Args[2]+".child")
			child.Stdout, child.Stderr = os.Stdout, os.Stderr
			if err := child.Start(); err != nil {
				panic(err)
			}
			deadline := time.Now().Add(5 * time.Second)
			for {
				if _, err := os.Stat(os.Args[2] + ".child"); err == nil {
					break
				}
				if time.Now().After(deadline) {
					panic("child not ready")
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
		if err := os.WriteFile(os.Args[2], []byte(fmt.Sprint(os.Getpid())), 0600); err != nil {
			panic(err)
		}
		if os.Args[1] == "orphan" {
			return
		}
		for {
			<-signals
			if os.Args[1] != "ignore" {
				_ = os.WriteFile(os.Args[2]+".stopped", []byte("graceful"), 0600)
				return
			}
		}
	default:
		os.Exit(2)
	}
}
