package app

import (
	"os"

	"golang.org/x/sys/windows"
)

func configSecurityDescriptor() (*windows.SECURITY_DESCRIPTOR, error) {
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, err
	}
	return windows.SecurityDescriptorFromString("D:P(A;;FA;;;SY)(A;;FA;;;" + user.User.Sid.String() + ")")
}

// Apply before writing any profile or recovery bytes. Do not inherit access
// from a custom --config directory that might be shared with other users.
func restrictConfigFile(f *os.File) error {
	sd, err := configSecurityDescriptor()
	if err != nil {
		return err
	}
	acl, _, err := sd.DACL()
	if err != nil {
		return err
	}
	// os.CreateTemp does not request WRITE_DAC; the named API opens the
	// security handle with the required rights while the empty file is held.
	return windows.SetNamedSecurityInfo(f.Name(), windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil)
}
