package persistence

import (
	"time"

	"github.com/google/uuid"
	"github.com/zelalem-12/distributed-file-storage/internal/domain"
)

type File struct {
	ID        string `gorm:"type:char(36);primary_key"`
	Name      string `gorm:"not null"`
	Path      string `gorm:"not null unique"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

func (modelFile *File) FromDomainModel(domainFile *domain.File) {

	modelFile.ID = domainFile.GetID().String()
	modelFile.Name = domainFile.GetName()
	modelFile.Path = domainFile.GetPath()
	modelFile.CreatedAt = domainFile.GetCreatedAt()
	modelFile.UpdatedAt = domainFile.GetUpdatedAt()
}

func (modelFile *File) ToDomainModel() *domain.File {
	domainFile := &domain.File{}

	fileID, _ := uuid.Parse(modelFile.ID)
	domainFile.SetID(fileID)
	domainFile.SetName(modelFile.Name)
	domainFile.SetPath(modelFile.Path)
	domainFile.SetCreatedAt(modelFile.CreatedAt)
	domainFile.SetUpdatedAt(modelFile.UpdatedAt)

	return domainFile
}
