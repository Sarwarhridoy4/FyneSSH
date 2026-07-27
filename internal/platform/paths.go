package platform

import (
	"os"
	"path/filepath"
)

func SSHBasePath() string {
	return filepath.Join(userHome(), ".ssh")
}

func AuthorizedKeysPath() string {
	return filepath.Join(SSHBasePath(), "authorized_keys")
}

func ConfigPath() string {
	return filepath.Join(SSHBasePath(), "config")
}

func PrivateKeyPath(name string) string {
	return filepath.Join(SSHBasePath(), name)
}

func PublicKeyPath(name string) string {
	return PrivateKeyPath(name) + ".pub"
}

func userHome() string {
	home, _, _ := os.UserHomeDir()
	return home
}
