package es

import (
	"emuSync/models"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/mholt/archiver/v3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	colorInfo    = color.New(color.FgCyan).SprintFunc()
	colorSuccess = color.New(color.FgGreen).SprintFunc()
	colorWarning = color.New(color.FgYellow).SprintFunc()
	colorError   = color.New(color.FgRed).SprintFunc()
	colorBold    = color.New(color.Bold).SprintFunc()
)

func (es *EmuSync) BackupDevice(device models.Device, backupROMs bool) error {
	fmt.Println(colorBold(fmt.Sprintf("Backing up device %s [%s]...", device.ID, device.Nickname)))

	deviceBackupDir := filepath.Join(backupDir, device.ID)

	sourceDirs := map[string]string{
		device.Directories.EmulatorConfigs: filepath.Join(deviceBackupDir, configsDirectory),
		device.Directories.Bios:            filepath.Join(deviceBackupDir, biosDirectory),
		device.Directories.Screenshots:     filepath.Join(deviceBackupDir, screenshotsDirectory),
		device.Directories.Overlays:        filepath.Join(deviceBackupDir, overlaysDirectory),
	}

	if backupROMs {
		sourceDirs[device.Directories.Roms] = filepath.Join(deviceBackupDir, romsDirectory)
	} else {
		fmt.Println(colorWarning("Skipping ROM backup!\n"))
	}

	err := es.BackupSaves(device)
	if err != nil {
		return err
	}

	err = es.BackupSaveStates(device)
	if err != nil {
		return err
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

		fmt.Printf("%s %s\n",
			colorInfo("Backing up:"),
			fmt.Sprintf("%s → %s", sourceDir, deviceBackupDir))

		if err := es.backupDirectory(device, sourceDir, deviceBackupDir); err != nil {
			fmt.Printf("%s %s: %v\n",
				colorWarning("Warning: failed to backup"),
				sourceDir, err)
			continue
		}
	}

	fmt.Printf("\n%s %s\n",
		colorSuccess("Backup completed for device"),
		colorBold(fmt.Sprintf("%s [%s]", device.ID, device.Nickname)))
	return nil
}

func (es *EmuSync) BackupSaves(device models.Device) error {
	return es.backupDirectoryWithArchive(device, device.Directories.Saves, savesDirectory)
}

func (es *EmuSync) BackupSaveStates(device models.Device) error {
	return es.backupDirectoryWithArchive(device, device.Directories.SaveStates, saveStatesDirectory)
}

func (es *EmuSync) backupDirectoryWithArchive(device models.Device, remotePath, backupType string) error {
	if remotePath == "" {
		return fmt.Errorf("no %s directory configured for device %s", backupType, device.ID)
	}

	fmt.Printf("%s %s for device %s [%s]\n",
		colorInfo("Creating timestamped backup of"),
		colorBold(backupType),
		colorBold(device.ID),
		device.Nickname)

	timestamp := time.Now().Format("2006-01-02T15-04-05")
	backupFilename := fmt.Sprintf("emusync_%s_%s_%s.zip", strings.ToLower(strings.ReplaceAll(backupType, " ", "_")), device.ID, timestamp)
	manifestFilename := fmt.Sprintf("emusync_%s_%s_%s.json", strings.ToLower(strings.ReplaceAll(backupType, " ", "_")), device.ID, timestamp)

	deviceBackupDir := filepath.Join(backupDir, device.ID)
	archiveDir := filepath.Join(deviceBackupDir, backupType)
	backupPath := filepath.Join(archiveDir, backupFilename)
	manifestPath := filepath.Join(archiveDir, manifestFilename)

	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	tempDir := filepath.Join(backupDir, fmt.Sprintf("temp_%s_%s_%s", strings.ToLower(strings.ReplaceAll(backupType, " ", "_")), device.ID, timestamp))
	defer os.RemoveAll(tempDir)

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	fmt.Printf("%s %s from device...\n",
		colorInfo("Pulling"),
		colorBold(backupType))
	if err := es.backupDirectoryContents(device, remotePath, tempDir); err != nil {
		return fmt.Errorf("failed to pull %s: %w", backupType, err)
	}

	manifest, err := es.createFileManifest(tempDir, device.ID, timestamp, backupFilename)
	if err != nil {
		return fmt.Errorf("failed to create file manifest: %w", err)
	}

	if err := es.writeManifestToFile(manifest, manifestPath); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	fmt.Printf("%s %s...\n",
		colorInfo("Creating ZIP archive"),
		colorBold(backupFilename))

	var filesToArchive []string
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	for _, entry := range entries {
		filesToArchive = append(filesToArchive, filepath.Join(tempDir, entry.Name()))
	}

	err = archiver.Archive(filesToArchive, backupPath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP archive: %w", err)
	}

	fmt.Printf("%s %s\n",
		colorSuccess("Manifest created:"),
		colorBold(manifestPath))
	fmt.Printf("%s %s\n\n",
		colorSuccess("Backup completed:"),
		colorBold(backupPath))
	return nil
}

func (es *EmuSync) backupDirectoryContents(device models.Device, remotePath, localPath string) error {
	files, err := es.ListFiles(device, remotePath)
	if err != nil {
		return fmt.Errorf("failed to list files in %s: %w", remotePath, err)
	}

	for _, file := range files {
		localFilePath := filepath.Join(localPath, file.Name)

		if file.IsDirectory {
			if err := es.backupDirectory(device, file.Path, localFilePath); err != nil {
				return fmt.Errorf("failed to backup directory %s: %w", file.Path, err)
			}
		} else {
			localFile, err := os.Create(localFilePath)
			if err != nil {
				return fmt.Errorf("failed to create local file %s: %w", localFilePath, err)
			}

			err = device.ADBDevice.Pull(file.Path, localFile)
			localFile.Close()

			if err != nil {
				return fmt.Errorf("failed to pull file %s: %w", file.Path, err)
			}
		}
	}

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

func (es *EmuSync) createFileManifest(rootDir string, deviceID string, timestamp string, archiveName string) (*models.BackupManifest, error) {
	saves := make(map[string][]interface{})

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			dirContents, err := es.getDirectoryContents(filepath.Join(rootDir, entry.Name()))
			if err != nil {
				continue
			}
			saves[entry.Name()] = dirContents
		} else {
			if saves["root"] == nil {
				saves["root"] = []interface{}{}
			}
			saves["root"] = append(saves["root"], entry.Name())
		}
	}

	manifest := &models.BackupManifest{
		DeviceID:    deviceID,
		Timestamp:   timestamp,
		ArchiveName: archiveName,
		Saves:       saves,
	}

	return manifest, nil
}

func (es *EmuSync) getDirectoryContents(dirPath string) ([]interface{}, error) {
	var contents []interface{}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return contents, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subContents, err := es.getDirectoryContents(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				continue
			}
			nestedMap := map[string][]interface{}{
				entry.Name(): subContents,
			}
			contents = append(contents, nestedMap)
		} else {
			contents = append(contents, entry.Name())
		}
	}

	return contents, nil
}

func (es *EmuSync) writeManifestToFile(manifest *models.BackupManifest, filePath string) error {
	jsonData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0644)
}
