package es

import (
	"emuSync/models"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var dataDir = filepath.Join(models.EmuSyncDataRoot, "config")

func (es *EmuSync) DoesConfigExist(id string) bool {
	configPath := filepath.Join(dataDir, fmt.Sprintf("%s.json", id))

	_, err := os.Stat(configPath)
	return err == nil
}

func (es *EmuSync) InitDevice(id string) (*models.Device, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return &models.Device{}, fmt.Errorf("failed to create data directory: %w", err)
	}

	configPath := filepath.Join(dataDir, fmt.Sprintf("%s.json", id))

	if es.DoesConfigExist(configPath) {
		return &models.Device{}, fmt.Errorf("config file already exists for device: %s", id)
	}

	device, err := es.GetDevice(id)
	if err != nil {
		return &models.Device{}, err
	}

	templateDevice := models.Device{
		ID:    device.ID,
		Model: device.Model,
		Directories: models.DeviceDirectoryPaths{
			EmulatorConfigs: "",
			Bios:            "",
			Roms:            "",
			Saves:           "",
			SaveStates:      "",
			Screenshots:     "",
		},
	}

	data, err := json.MarshalIndent(templateDevice, "", "  ")
	if err != nil {
		return &models.Device{}, fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return &models.Device{}, fmt.Errorf("failed to write config file: %w", err)
	}

	return &templateDevice, nil
}

func (es *EmuSync) LoadConfig(id string) (*models.Device, error) {
	configPath := filepath.Join(dataDir, fmt.Sprintf("%s.json", id))

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found for device ID: %s", id)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var device models.Device
	if err := json.Unmarshal(data, &device); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &device, nil
}
