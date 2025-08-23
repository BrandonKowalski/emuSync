package es

import (
	"emuSync/models"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var backupDir = filepath.Join(models.EmuSyncDataRoot, "backups")

func (es *EmuSync) BackupDevice(device models.Device, backupROMs bool) error {
	fmt.Println(fmt.Sprintf("Backing up device %s [%s]...", device.ID, device.Nickname))

	deviceBackupDir := filepath.Join(backupDir, device.ID)

	sourceDirs := map[string]string{
		device.Directories.EmulatorConfigs: filepath.Join(deviceBackupDir, filepath.Base(device.Directories.EmulatorConfigs)),
		device.Directories.Bios:            filepath.Join(deviceBackupDir, filepath.Base(device.Directories.Bios)),
		device.Directories.Saves:           filepath.Join(deviceBackupDir, filepath.Base(device.Directories.Saves)),
		device.Directories.SaveStates:      filepath.Join(deviceBackupDir, filepath.Base(device.Directories.SaveStates)),
		device.Directories.Screenshots:     filepath.Join(deviceBackupDir, filepath.Base(device.Directories.Screenshots)),
		device.Directories.Overlays:        filepath.Join(deviceBackupDir, filepath.Base(device.Directories.Overlays)),
	}

	if backupROMs {
		sourceDirs[device.Directories.Roms] = filepath.Join(deviceBackupDir, filepath.Base(device.Directories.Roms))
	} else {
		fmt.Println("Skipping ROM backup")
	}

	for _, dest := range sourceDirs {
		if err := os.MkdirAll(dest, 0755); err != nil {
			return fmt.Errorf("failed to create backup directory %s: %w", dest, err)
		}
	}

	for sourceDir, _ := range sourceDirs {
		if sourceDir == "" {
			continue
		}

		fmt.Printf("Backing up %s to %s...\n", sourceDir, deviceBackupDir)

		if err := es.backupDirectory(device, sourceDir, deviceBackupDir); err != nil {
			fmt.Printf("Warning: failed to backup %s: %v\n", sourceDir, err)
			continue
		}
	}

	fmt.Printf("Backup completed for device %s\n", device.ID)
	return nil
}

func (es *EmuSync) backupDirectory(device models.Device, sourceDir, destDir string) error {
	cmd := exec.Command("adb", "-s", device.ID, "pull", sourceDir, destDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to pull directory %s to %s: %w\nOutput: %s", sourceDir, destDir, err, string(output))
	}

	return nil
}
