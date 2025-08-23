package models

import "github.com/electricbubble/gadb"

type Device struct {
	ID          string               `json:"id"`
	Nickname    string               `json:"nickname"`
	Model       string               `json:"model"`
	ADBDevice   *gadb.Device         `json:"-"`
	Directories DeviceDirectoryPaths `json:"directories"`
}

type DeviceDirectoryPaths struct {
	EmulatorConfigs string `json:"emulator_configs"`
	Bios            string `json:"bios"`
	Roms            string `json:"roms"`
	Saves           string `json:"saves"`
	SaveStates      string `json:"save_states"`
	Overlays        string `json:"overlays"`
	Screenshots     string `json:"screenshots"`
}
