package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/zelalem-12/distributed-file-storage/internal/domain"
	"gorm.io/gorm"
)

type fileRepo struct {
	db *gorm.DB
}

func NewFileRepo(db *gorm.DB) domain.FileRepo {
	return fileRepo{
		db,
	}
}

func (repo fileRepo) Insert(ctx context.Context, entity *domain.File) (*domain.File, error) {
	model := &File{}
	model.FromDomainModel(entity)

	err := repo.
		db.
		WithContext(ctx).
		Create(model).
		Error

	if err != nil {
		return nil, err
	}

	return model.ToDomainModel(), nil
}
func (repo fileRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	model := File{}

	err := repo.
		db.WithContext(ctx).
		Where("id = ?", id).
		First(&model).
		Error

	if err != nil {
		return nil, err
	}

	return model.ToDomainModel(), nil
}
func (repo fileRepo) FindAll(ctx context.Context) ([]*domain.File, error) {
	models := []*File{}

	err := repo.
		db.
		WithContext(ctx).
		Find(&models).
		Error

	if err != nil {
		return []*domain.File{}, err
	}

	var domians []*domain.File
	for _, model := range models {
		domians = append(domians, model.ToDomainModel())
	}

	return domians, nil
}
