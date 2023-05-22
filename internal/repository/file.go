package repository

import (
	"github.com/google/uuid"
	"github.com/taranovegor/uploadly/internal/domain"
	"gorm.io/gorm"
)

type FileRepository struct {
	domain.InterfaceFileRepository
	orm *gorm.DB
}

func NewFileRepository(
	orm *gorm.DB,
) FileRepository {
	return FileRepository{
		orm: orm,
	}
}

func (repository FileRepository) Save(file domain.File) error {
	return repository.orm.Create(file).Error
}

func (repository FileRepository) Get(id uuid.UUID) (domain.File, error) {
	var file domain.File
	if err := repository.orm.First(&file, id).Error; err != nil {
		return file, err
	}

	return file, nil
}

func (repository FileRepository) Delete(id uuid.UUID) error {
	return repository.orm.Delete(&domain.File{}, id).Error
}
