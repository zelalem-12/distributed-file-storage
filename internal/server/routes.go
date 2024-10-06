package server

import (
	"github.com/julienschmidt/httprouter"
)

type Group struct {
	Prefix string
	Router *httprouter.Router
}

// NewGroup creates a new group with the specified prefix.
func NewGroup(prefix string, router *httprouter.Router) *Group {
	return &Group{Prefix: prefix, Router: router}
}

// GET adds a new GET route to the group.
func (g *Group) GET(path string, handle httprouter.Handle) {
	g.Router.GET(g.Prefix+path, handle)
}

// POST adds a new POST route to the group.
func (g *Group) POST(path string, handle httprouter.Handle) {
	g.Router.POST(g.Prefix+path, handle)
}

func RegisterRoutes(handler *FileHandler) *Group {
	router := httprouter.New()
	// Create route groups
	apiRouter := NewGroup("/api/v1", router)

	apiRouter.Router.GET("/", handler.HomeHandler)
	apiRouter.POST("/upload", handler.UploadHandler)
	apiRouter.GET("/downloads", handler.DownloadsHandler)
	apiRouter.GET("/downloads/:id", handler.GetFileHandler)

	return apiRouter
}
