package storage

import (
	"github.com/taranovegor/uploadly/internal/domain"
	"io"
)

type InterfaceStorage interface {
	Save(domain.File, io.Reader) error
	Url(domain.File) (string, error)
	Delete(domain.File) error
}
