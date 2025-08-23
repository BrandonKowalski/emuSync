package es

import (
	"emuSync/models"
	"errors"
	"github.com/electricbubble/gadb"
	"slices"
	"strings"
)

func (es *EmuSync) ListDevices() ([]models.Device, error) {
	dl, err := client.DeviceList()
	if err != nil {
		return nil, err
	}

	var devices []models.Device

	for _, device := range dl {
		parsed, err := parseDevice(device)
		if err != nil {
			continue
		}

		devices = append(devices, parsed)
	}

	slices.SortFunc(devices, func(a, b models.Device) int {
		return strings.Compare(a.Model, b.Model)
	})

	return devices, err
}

func (es *EmuSync) GetDeviceWithConfig(id string) (models.Device, error) {
	d, err := es.GetDevice(id)
	if err != nil {
		return models.Device{}, err
	}

	config, err := es.LoadConfig(id)
	if err != nil {
		return models.Device{}, err
	}

	config.ADBDevice = d.ADBDevice

	return *config, nil
}

func (es *EmuSync) GetDevice(id string) (models.Device, error) {
	dl, err := client.DeviceList()
	if err != nil {
		return models.Device{}, err
	}

	if len(dl) == 0 {
		return models.Device{}, errors.New("No devices found")
	}

	for _, d := range dl {
		if d.Serial() == id {
			parsed, err := parseDevice(d)
			if err != nil {
				return models.Device{}, err
			}

			return parsed, nil
		}
	}

	return models.Device{}, errors.New("device not found")
}

func parseDevice(device gadb.Device) (models.Device, error) {
	model, err := device.Model()
	if err != nil {
		return models.Device{}, err
	}

	model = strings.ReplaceAll(model, "_", " ")
	model = strings.ReplaceAll(model, "-", " ")
	model = strings.TrimSpace(model)

	return models.Device{
		ID:        device.Serial(),
		Model:     model,
		ADBDevice: &device,
	}, nil
}
