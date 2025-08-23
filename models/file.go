package models

import "time"

type File struct {
	Name         string    `json:"name,omitempty"`
	Path         string    `json:"path,omitempty"`
	IsDir        bool      `json:"is_dir,omitempty"`
	Size         uint32    `json:"size,omitempty"`
	LastModified time.Time `json:"last_modified"`
}
