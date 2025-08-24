package es

import (
	"path/filepath"
)

const (
	emuSyncRoot          = "/Users/btk/emuSync"
	configsDirectory     = "Configs"
	biosDirectory        = "BIOS"
	romsDirectory        = "Roms"
	savesDirectory       = "Saves"
	saveStatesDirectory  = "Save States"
	screenshotsDirectory = "Screenshots"
	overlaysDirectory    = "Overlays"
)

var backupDir = filepath.Join(emuSyncRoot, "backups")
