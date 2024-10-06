package persistence

import (
	"log"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {

	fileEntities := []interface{}{
		&File{},
	}

	db.AutoMigrate(fileEntities...)

	log.Println("DB Schema Migrated...")
}
