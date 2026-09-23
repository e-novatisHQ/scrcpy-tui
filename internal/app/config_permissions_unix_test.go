//go:build !windows

package app

import (
	"os"
	"testing"
)

func assertPrivateConfig(t *testing.T, path string) {
	t.Helper()
	st, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if st.Mode().Perm() != 0600 {
		t.Fatalf("profile permissions: %v", st.Mode())
	}
}
