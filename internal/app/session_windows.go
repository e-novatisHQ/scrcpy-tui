package app

import (
	"errors"
	"os"
	"os/exec"
)

// Fail closed until the Windows ownership backend is implemented (WIN-02).
func runSession(_ *exec.Cmd, _ <-chan os.Signal) error {
	return errors.New("sessions Windows indisponibles : backend de processus non implémenté")
}
