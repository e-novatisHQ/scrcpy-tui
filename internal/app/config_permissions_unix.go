//go:build !windows

package app

import "os"

func restrictConfigFile(f *os.File) error { return f.Chmod(0600) }
