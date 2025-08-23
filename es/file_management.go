package es

import (
	"emuSync/models"
	"path/filepath"
	"slices"
	"strings"
)

func (es *EmuSync) ListFiles(device models.Device, path string) ([]models.File, error) {
	adb := device.ADBDevice

	dfi, err := adb.List(path)
	if err != nil {
		return nil, err
	}

	var files []models.File

	for _, file := range dfi {
		if strings.HasPrefix(file.Name, ".") {
			continue
		}

		files = append(files, models.File{
			Name:         file.Name,
			Path:         filepath.Join(path, file.Name),
			IsDirectory:  file.Mode == 16888, // TODO determine if this is correct... the library is not detecting directories properly
			Size:         file.Size,
			LastModified: file.LastModified,
		})
	}

	slices.SortFunc(files, func(a, b models.File) int {
		return strings.Compare(a.Name, b.Name)
	})

	return files, nil
}
