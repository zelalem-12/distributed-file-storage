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
	_type     string
	size      int64
	extension string
	createdAt *time.Time
}

func CreateFile(name, path, _type, extension string, size int64) (*File, error) {

	now := time.Now()
	fileID := uuid.New()

	file := &File{}
	file.SetID(fileID)
	if err := file.SetName(name); err != nil {
		return nil, err
	}
	if err := file.SetPath(path); err != nil {
		return nil, err
	}
	file.SetType(_type)
	file.SetSize(size)
	file.SetExtension(extension)
	file.SetCreatedAt(&now)

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
func (file *File) SetType(_type string) {
	file._type = _type
}
func (file *File) SetSize(size int64) {
	file.size = size
}
func (file *File) SetExtension(extension string) {
	file.extension = extension
}

func (file *File) SetCreatedAt(createdAt *time.Time) {
	file.createdAt = createdAt
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
func (file *File) GetType() string {
	return file._type
}
func (file *File) GetSize() int64 {
	return file.size
}
func (file *File) GetExtension() string {
	return file.extension
}

func (file *File) GetCreatedAt() *time.Time {
	return file.createdAt
}
