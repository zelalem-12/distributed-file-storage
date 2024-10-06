package server

import "github.com/google/uuid"

type Request struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Path string    `json:"path"`
}
