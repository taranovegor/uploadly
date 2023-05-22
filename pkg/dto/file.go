package dto

import (
	"github.com/google/uuid"
	"time"
)

type FileSaveSlice []interface{}
type FilesSaveMap map[string]FileSaveSlice

type File struct {
	Id             uuid.UUID `json:"id"`
	Context        string    `json:"context"`
	OriginFilename string    `json:"origin_filename"`
	MimeType       string    `json:"mime_type"`
	Size           uint      `json:"size"`
	CreatedAt      time.Time `json:"created_at"`
	Url            string    `json:"url"`
}
