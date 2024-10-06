package main

import (
	"log"

	"github.com/zelalem-12/distributed-file-storage/internal/persistence"
	"github.com/zelalem-12/distributed-file-storage/internal/server"
	"github.com/zelalem-12/distributed-file-storage/internal/services"
	"github.com/zelalem-12/distributed-file-storage/pkg"
	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		fx.Provide(
			pkg.Load,
			pkg.InitPostgresDB,
			persistence.NewFileRepo,
			services.NewFileService,
			server.NewFileHandler,
			server.RegisterRoutes,
		),
		//	fx.Invoke(fx.Invoke(EchoSwaggerInit)),
		fx.Invoke(persistence.MigrateDB),
		fx.Invoke(pkg.StartHTTPServer),
	)

	log.Println("Starting the application...")
	app.Run() // This runs the application and manages lifecycle

}

// func EchoSwaggerInit(e *echo.Echo) {
// 	e.GET("/swagger/*", echoSwagger.WrapHandler)
// }
