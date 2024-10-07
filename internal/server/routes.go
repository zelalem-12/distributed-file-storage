package server

import (
	"github.com/julienschmidt/httprouter"
)

type Group struct {
	Prefix string
	Router *httprouter.Router
}

func NewGroup(prefix string, router *httprouter.Router) *Group {
	return &Group{Prefix: prefix, Router: router}
}

func (g *Group) GET(path string, handle httprouter.Handle) {
	g.Router.GET(g.Prefix+path, handle)
}

func (g *Group) POST(path string, handle httprouter.Handle) {
	g.Router.POST(g.Prefix+path, handle)
}

func RegisterRoutes(handler *FileHandler) *Group {
	router := httprouter.New()
	apiRouter := NewGroup("/api/v1", router)

	apiRouter.Router.GET("/", handler.HomeHandler)
	apiRouter.POST("/upload", handler.UploadFilesHandler)
	apiRouter.GET("/files", handler.GetFilesDataHandler)
	apiRouter.GET("/downloads/:id", handler.DownloadFileByIDHandler)

	return apiRouter
}
