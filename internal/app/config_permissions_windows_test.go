package app

import (
	"testing"

	"golang.org/x/sys/windows"
)

func assertPrivateConfig(t *testing.T, path string) {
	t.Helper()
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := configSecurityDescriptor()
	if err != nil {
		t.Fatal(err)
	}
	if sd.String() != expected.String() {
		t.Fatalf("profile ACL: %s; want %s", sd.String(), expected.String())
	}
}
