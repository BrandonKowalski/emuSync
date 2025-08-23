package models

type Config struct {
	DevicePaths DeviceDirectoryPaths `json:"device_paths"`
}

type DeviceDirectoryPaths struct {
	EmulatorConfigs string `json:"emulator_configs,omitempty"`
	Bios            string `json:"bios,omitempty"`
	Roms            string `json:"roms,omitempty"`
	Saves           string `json:"saves,omitempty"`
	SaveStates      string `json:"save_states,omitempty"`
	Screenshots     string `json:"screenshots,omitempty"`
}
