package platform

import (
	"os"
	"runtime"
)

// GetDefaultShell returns the user's default shell, or a sensible fallback.
func GetDefaultShell() string {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}
	if runtime.GOOS == "windows" {
		if comspec := os.Getenv("ComSpec"); comspec != "" {
			return comspec
		}
		return "powershell.exe"
	}
	return "/bin/sh"
}
