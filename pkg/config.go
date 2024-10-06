package pkg

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	SERVER_PORT       int    `mapstructure:"SERVER_PORT"  validate:"required"`
	POSTGRES_HOST     string `mapstructure:"POSTGRES_HOST" validate:"required"`
	POSTGRES_PORT     int    `mapstructure:"POSTGRES_PORT" validate:"required"`
	POSTGRES_DATABASE string `mapstructure:"POSTGRES_DATABASE" validate:"required"`
	POSTGRES_USER     string `mapstructure:"POSTGRES_USER" validate:"required"`
	POSTGRES_PASSWORD string `mapstructure:"POSTGRES_PASSWORD" validate:"required"`
}

func Load() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal("Server port Required")
	}

	if os.Getenv("POSTGRES_HOST") == "" {
		log.Fatal("Postgres host is  required")
	}

	postGresPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Fatal("Postgres port required")
	}

	if os.Getenv("POSTGRES_DATABASE") == "" {
		log.Fatal("Postgres database is  required")
	}

	return Config{
		SERVER_PORT:       serverPort,
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     postGresPort,
		POSTGRES_DATABASE: os.Getenv("POSTGRES_DATABASE"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
	}, nil

}
