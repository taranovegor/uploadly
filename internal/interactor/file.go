package interactor

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/taranovegor/uploadly/internal/builder"
	"github.com/taranovegor/uploadly/internal/domain"
	"github.com/taranovegor/uploadly/internal/storage"
	"github.com/taranovegor/uploadly/pkg/config"
	"github.com/taranovegor/uploadly/pkg/dto"
	"mime/multipart"
	"net/http"
)

type FileInteractor struct {
	domain.InterfaceFileInteractor
	contexts   config.FileContexts
	repository domain.InterfaceFileRepository
	builder    builder.File
	storage    storage.InterfaceStorage
}

func NewFileInteractor(
	contexts config.FileContexts,
	repository domain.InterfaceFileRepository,
	builder builder.File,
	storage storage.InterfaceStorage,
) FileInteractor {
	return FileInteractor{
		contexts:   contexts,
		repository: repository,
		builder:    builder,
		storage:    storage,
	}
}

func (interactor FileInteractor) Save(contextName string, filename string, file multipart.File) (*dto.File, error) {
	context, err := interactor.contexts.GetContext(contextName)
	if err != nil {
		return nil, err
	}

	fileBuffer := new(bytes.Buffer)
	if _, err := fileBuffer.ReadFrom(file); err != nil {
		return nil, err
	}

	mimeType := http.DetectContentType(fileBuffer.Bytes())
	if !context.IsMimeTypeAllowed(mimeType) {
		return nil, errors.New(fmt.Sprintf("Mime-Type `%s` is not allowed for `%s` context is not allowed for this context", mimeType, contextName))
	}

	fileSize := fileBuffer.Len()
	if !context.IsFileSizeAcceptable(fileSize) {
		return nil, errors.New(fmt.Sprintf("File size is `%dB` and over `%dB` limit", fileSize, context.MaxFileSize))
	}

	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}

	entity := domain.NewFile(contextName, filename, mimeType, uint(fileSize))
	if err := interactor.repository.Save(entity); err != nil {
		return nil, err
	}

	if err := interactor.storage.Save(entity, file); err != nil {
		_ = interactor.repository.Delete(entity.Id)

		return nil, err
	}

	return interactor.builder.Build(entity)
}

func (interactor FileInteractor) Get(id uuid.UUID) (*dto.File, error) {
	entity, err := interactor.repository.Get(id)
	if err != nil {
		return nil, err
	}

	return interactor.builder.Build(entity)
}

func (interactor FileInteractor) Delete(id uuid.UUID) error {
	entity, err := interactor.repository.Get(id)
	if err != nil {
		return err
	}

	if err := interactor.storage.Delete(entity); err != nil {
		return err
	}

	return interactor.repository.Delete(id)
}
