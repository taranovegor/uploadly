package builder

import (
	"github.com/taranovegor/uploadly/internal/domain"
	"github.com/taranovegor/uploadly/internal/storage"
	"github.com/taranovegor/uploadly/pkg/dto"
)

type File struct {
	storage storage.InterfaceStorage
}

func NewFile(
	storage storage.InterfaceStorage,
) File {
	return File{
		storage: storage,
	}
}

func (builder File) Build(file domain.File) (*dto.File, error) {
	url, err := builder.storage.Url(file)
	if err != nil {
		return nil, err
	}

	return &dto.File{
		Id:             file.Id,
		Context:        file.Context,
		OriginFilename: file.OriginFilename,
		MimeType:       file.MimeType,
		Size:           file.Size,
		CreatedAt:      file.CreatedAt,
		Url:            url,
	}, nil
}
