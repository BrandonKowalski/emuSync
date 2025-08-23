package es

import (
	"emuSync/models"
	"errors"
	"github.com/electricbubble/gadb"
	"slices"
	"strings"
)

func (s EmuSync) ListDevices() ([]models.Device, error) {
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

func (s EmuSync) GetDevice(id string) (gadb.Device, error) {
	dl, err := client.DeviceList()
	if err != nil {
		return gadb.Device{}, err
	}

	if len(dl) == 0 {
		return gadb.Device{}, errors.New("no devices found")
	}

	for _, d := range dl {
		if d.Serial() == id {
			return d, nil
		}
	}

	return gadb.Device{}, errors.New("device not found")
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
		ID:    device.Serial(),
		Model: model,
	}, nil
}
