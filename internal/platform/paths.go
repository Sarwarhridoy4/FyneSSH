package platform

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

// SSHBasePath returns the path to the user's .ssh directory.
func SSHBasePath() string {
	return filepath.Join(userHome(), ".ssh")
}

// AuthorizedKeysPath returns the expected path for authorized_keys.
func AuthorizedKeysPath() string {
	return filepath.Join(SSHBasePath(), "authorized_keys")
}

// ConfigPath returns the expected path for the SSH config file.
func ConfigPath() string {
	return filepath.Join(SSHBasePath(), "config")
}

// PrivateKeyPath returns the filesystem path for a named private key.
func PrivateKeyPath(name string) string {
	return filepath.Join(SSHBasePath(), name)
}

// PublicKeyPath returns the filesystem path for a named public key.
func PublicKeyPath(name string) string {
	return PrivateKeyPath(name) + ".pub"
}

func userHome() string {
	home, _ := os.UserHomeDir()
	if home != "" {
		return home
	}
	if runtime.GOOS == "windows" {
		if profile := os.Getenv("USERPROFILE"); profile != "" {
			return profile
		}
	}
	return "."
}

// CurrentUser returns the current OS user.
func CurrentUser() string {
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	if runtime.GOOS == "windows" {
		if username := os.Getenv("USERNAME"); username != "" {
			return username
		}
	}
	return "user"
}

// IsRootUser reports whether the current effective user is root/admin.
func IsRootUser() bool {
	if runtime.GOOS == "windows" {
		_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
		return err == nil
	}
	return os.Geteuid() == 0
}
