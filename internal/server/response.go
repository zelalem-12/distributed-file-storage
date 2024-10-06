package server

import "github.com/google/uuid"

type Response struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt string    `json:"created_at"`
}
