package app

import (
	"strings"
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
	control, _, err := sd.Control()
	if err != nil || control&windows.SE_DACL_PROTECTED == 0 {
		t.Fatalf("ACL inheritance is not protected: %v", err)
	}
	// Windows adds the AI bookkeeping flag even to a protected DACL. Compare
	// the actual ACE list, while checking protection independently above.
	actualACEs := sd.String()[strings.Index(sd.String(), "("):]
	expectedACEs := expected.String()[strings.Index(expected.String(), "("):]
	if actualACEs != expectedACEs {
		t.Fatalf("profile ACL: %s; want %s", sd.String(), expected.String())
	}
}
