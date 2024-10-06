package domain

import (
	"context"

	"github.com/google/uuid"
)

type FileRepo interface {
	Insert(ctx context.Context, entity *File) (*File, error)
	FindByID(ctx context.Context, id uuid.UUID) (*File, error)
	FindAll(ctx context.Context) ([]*File, error)
}
