package domain

import (
	"github.com/google/uuid"
	"github.com/taranovegor/uploadly/pkg/dto"
	"mime/multipart"
	"time"
)

type File struct {
	Id             uuid.UUID `gorm:"primary_key;size:36;<-:create"`
	Context        string    `gorm:"size:255"`
	OriginFilename string    `gorm:"size:255"`
	MimeType       string    `gorm:"size:128"`
	Size           uint
	CreatedAt      time.Time
}

type InterfaceFileRepository interface {
	Save(File) error
	Get(uuid.UUID) (File, error)
	Delete(uuid.UUID) error
}

type InterfaceFileInteractor interface {
	Save(contextName string, filename string, file multipart.File) (*dto.File, error)
	Get(id uuid.UUID) (*dto.File, error)
	Delete(id uuid.UUID) error
}

func NewFile(context string, filename string, mimeType string, size uint) File {
	return File{
		Id:             uuid.New(),
		Context:        context,
		OriginFilename: filename,
		MimeType:       mimeType,
		Size:           size,
		CreatedAt:      time.Now(),
	}
}
