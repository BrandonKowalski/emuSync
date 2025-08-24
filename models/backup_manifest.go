package models

type BackupManifest struct {
	DeviceID    string                   `json:"deviceId"`
	Timestamp   string                   `json:"timestamp"`
	ArchiveName string                   `json:"archiveName"`
	Saves       map[string][]interface{} `json:"saves"`
}
