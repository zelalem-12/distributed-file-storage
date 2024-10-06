package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/zelalem-12/distributed-file-storage/internal/domain"
)

type FileService struct {
	fileRepo domain.FileRepo
}

func NewFileService(fileRepo domain.FileRepo) FileServiceInterface {
	return &FileService{
		fileRepo: fileRepo,
	}
}

func (service FileService) Create(ctx context.Context, entity *domain.File) (*domain.File, error) {

	savedData, err := service.fileRepo.Insert(ctx, entity)
	if err != nil {
		return nil, err
	}
	return savedData, err
}
func (service FileService) GetById(ctx context.Context, id uuid.UUID) (*domain.File, error) {

	file, err := service.fileRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return file, err

}
func (service FileService) GetAll(ctx context.Context) ([]*domain.File, error) {

	files, err := service.fileRepo.FindAll(ctx)

	if err != nil {
		return nil, err
	}
	return files, err
}
