package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/zelalem-12/distributed-file-storage/internal/domain"
)

type FileServiceInterface interface {
	Create(ctx context.Context, entity *domain.File) (*domain.File, error)
	GetById(ctx context.Context, id uuid.UUID) (*domain.File, error)
	GetAll(ctx context.Context) ([]*domain.File, error)
}
