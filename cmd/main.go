package main

import (
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/zelalem-12/distributed-file-storage/docs"
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
		fx.Invoke(SwaggerInit),
	)

	log.Println("Starting the application...")
	app.Run() // This runs the application and manages lifecycle
}

func SwaggerInit(lc fx.Lifecycle, routerGroup *server.Group) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			routerGroup.Router.GET("/swagger/:any", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				httpSwagger.WrapHandler(w, r)
			})
			return nil
		},
	})
}
