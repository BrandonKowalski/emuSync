package models

import "time"

type File struct {
	Name         string    `json:"name,omitempty"`
	Path         string    `json:"path,omitempty"`
	IsDirectory  bool      `json:"is_directory,omitempty"`
	Size         uint32    `json:"size,omitempty"`
	LastModified time.Time `json:"last_modified"`
}
