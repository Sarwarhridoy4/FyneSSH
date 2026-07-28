package platform

import (
	"os"
	"runtime"
)

// DistroFamily represents a Linux distribution family.
type DistroFamily int

const (
	DistroUnknown DistroFamily = iota
	DistroDebian
	DistroRedHat
	DistroMacOS
)

func DetectDistro() DistroFamily {
	if runtime.GOOS == "darwin" {
		return DistroMacOS
	}
	_, err := os.Stat("/etc/debian_version")
	if err == nil {
		return DistroDebian
	}
	_, err = os.Stat("/etc/redhat-release")
	if err == nil {
		return DistroRedHat
	}
	return DistroUnknown
}
