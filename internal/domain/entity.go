package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type File struct {
	id        uuid.UUID
	name      string
	path      string
	createdAt *time.Time
	updatedAt *time.Time
	deletedAt *time.Time
}

func CreateFile(name, path string) (*File, error) {

	now := time.Now()
	fileID := uuid.New()

	file := &File{}
	file.SetID(fileID)
	if err := file.SetName(name); err != nil {
		return nil, err
	}
	if err := file.SetPath(name); err != nil {
		return nil, err
	}

	file.SetCreatedAt(&now)
	file.SetUpdatedAt(&now)

	return file, nil
}
func (file *File) SetID(id uuid.UUID) {
	file.id = id
}

func (file *File) SetName(name string) error {
	if len(name) == 0 {
		return errors.New("name is required")
	}
	file.name = name
	return nil
}

func (file *File) SetPath(path string) error {
	if len(path) == 0 {
		return errors.New("path is required")
	}
	file.path = path
	return nil
}

func (file *File) SetCreatedAt(createdAt *time.Time) {
	file.createdAt = createdAt
}

func (file *File) SetUpdatedAt(updatedAt *time.Time) {
	file.updatedAt = updatedAt
}

func (file *File) SetDeletedAt(deletedAt *time.Time) {
	file.deletedAt = deletedAt
}

func (file *File) GetID() uuid.UUID {
	return file.id
}

func (file *File) GetName() string {
	return file.name
}

func (file *File) GetPath() string {
	return file.path
}
func (file *File) GetCreatedAt() *time.Time {
	return file.createdAt
}

func (file *File) GetUpdatedAt() *time.Time {
	return file.updatedAt
}

func (file *File) GetDeletedAt() *time.Time {
	return file.deletedAt
}
