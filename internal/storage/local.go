package storage

import (
	"errors"
	"fmt"
	"github.com/taranovegor/uploadly/internal/config"
	"github.com/taranovegor/uploadly/internal/domain"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	InterfaceStorage
	rootPath string
}

func NewLocal(
	rootPath string,
) *Local {
	return &Local{
		rootPath: rootPath,
	}
}

func (storage Local) Save(file domain.File, src io.Reader) error {
	contextPath := storage.contextPath(file.Context)
	if _, err := os.Stat(contextPath); os.IsNotExist(err) {
		if err := os.Mkdir(contextPath, 0770); err != nil {
			return err
		}
	}

	f2, err := os.Create(storage.filePath(file))
	if err != nil {
		fmt.Println(err)
	}
	_, err = io.Copy(f2, src)

	return err
}

func (storage Local) Url(file domain.File) (string, error) {
	if !storage.isFileExists(file) {
		return "", errors.New(fmt.Sprintf("File `%s` was not found in context `%s`", file.Id.String(), file.Context))
	}

	return fmt.Sprintf("/%s/%s/%s", config.StaticHttpPath, file.Context, storage.filename(file)), nil
}

func (storage Local) Delete(file domain.File) error {
	if !storage.isFileExists(file) {
		return errors.New(fmt.Sprintf("File `%s` was not found in context `%s`", file.Id.String(), file.Context))
	}

	return os.Remove(storage.filePath(file))
}

func (storage Local) contextPath(context string) string {
	return fmt.Sprintf("%s/%s", storage.rootPath, context)
}

func (storage Local) filename(file domain.File) string {
	return file.Id.String() + filepath.Ext(file.OriginFilename)
}

func (storage Local) filePath(file domain.File) string {
	return fmt.Sprintf("%s/%s", storage.contextPath(file.Context), storage.filename(file))
}

func (storage Local) isFileExists(file domain.File) bool {
	_, err := os.Stat(storage.filePath(file))

	return err == nil
}
