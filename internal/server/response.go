package server

import "github.com/google/uuid"

type Response struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Type      string    `json:"type"`
	Size      int64     `json:"size"`
	Extension string    `json:"extension"`
	CreatedAt string    `json:"created_at"`
}
