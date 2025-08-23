package es

import (
	"emuSync/models"
	"github.com/electricbubble/gadb"
	"path/filepath"
	"slices"
	"strings"
)

func (s EmuSync) ListFiles(device gadb.Device, path string) ([]models.File, error) {
	dfi, err := device.List(path)
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
			IsDir:        file.Mode.IsDir(),
			Size:         file.Size,
			LastModified: file.LastModified,
		})
	}

	slices.SortFunc(files, func(a, b models.File) int {
		return strings.Compare(a.Name, b.Name)
	})

	return files, nil
}
