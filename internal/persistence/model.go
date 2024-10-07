package persistence

import (
	"time"

	"github.com/google/uuid"
	"github.com/zelalem-12/distributed-file-storage/internal/domain"
)

// File Name: The name of the file (e.g., document.pdf).
// File Path: If the file is stored on disk, store the path to the file (e.g., /uploads/document.pdf).
// File Type/MIME Type: The type of file (e.g., application/pdf, image/jpeg).
// File Size: The size of the file in bytes.
// Upload Date/Time: The timestamp when the file was uploaded.
// Uploader ID: If your system involves users, store who uploaded the file (could be a user ID or reference).
// File Extension: The extension of the file (e.g., .pdf, .jpg).

type File struct {
	ID        string `gorm:"type:char(36);primary_key"`
	Name      string `gorm:"not null"`
	Path      string `gorm:"not null unique"`
	Type      string
	Size      int64 `gorm:"type:int"`
	Extension string
	CreatedAt *time.Time
}

func (modelFile *File) FromDomainModel(domainFile *domain.File) {

	modelFile.ID = domainFile.GetID().String()
	modelFile.Name = domainFile.GetName()
	modelFile.Path = domainFile.GetPath()
	modelFile.Type = domainFile.GetType()
	modelFile.Size = domainFile.GetSize()
	modelFile.Extension = domainFile.GetExtension()
	modelFile.CreatedAt = domainFile.GetCreatedAt()
}

func (modelFile *File) ToDomainModel() *domain.File {
	domainFile := &domain.File{}

	fileID, _ := uuid.Parse(modelFile.ID)
	domainFile.SetID(fileID)
	domainFile.SetName(modelFile.Name)
	domainFile.SetPath(modelFile.Path)
	domainFile.SetType(modelFile.Type)
	domainFile.SetSize(modelFile.Size)
	domainFile.SetExtension(modelFile.Extension)
	domainFile.SetCreatedAt(modelFile.CreatedAt)

	return domainFile
}
