package updater

import (
	"testing"
)

func TestNewUpdater(t *testing.T) {
	version := "v1.0.0"
	verbose := true

	updater := NewUpdater(version, verbose)

	if updater.currentVer != version {
		t.Errorf("NewUpdater() currentVer = %v, want %v", updater.currentVer, version)
	}

	if updater.verbose != verbose {
		t.Errorf("NewUpdater() verbose = %v, want %v", updater.verbose, verbose)
	}
}
