package pkg

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresDB(cfg Config) (*gorm.DB, error) {

	host := cfg.POSTGRES_HOST
	port := cfg.POSTGRES_PORT
	user := cfg.POSTGRES_USER
	password := cfg.POSTGRES_PASSWORD
	dbName := cfg.POSTGRES_DATABASE

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connectionString,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
