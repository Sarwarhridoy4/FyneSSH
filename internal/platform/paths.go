package platform

import (
	"os"
	"path/filepath"
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
	return home
}
