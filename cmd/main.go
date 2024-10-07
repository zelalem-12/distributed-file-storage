package main

import (
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/zelalem-12/distributed-file-storage/docs/openapi" // Swagger docs
	"github.com/zelalem-12/distributed-file-storage/internal/persistence"
	"github.com/zelalem-12/distributed-file-storage/internal/server"
	"github.com/zelalem-12/distributed-file-storage/internal/services"
	"github.com/zelalem-12/distributed-file-storage/pkg"
	"go.uber.org/fx"
)

// @title Distributed File Storage API
// @version 1.0
// @description This is the API for a distributed file storage system.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url http://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1/

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

// SwaggerInit initializes the Swagger route
func SwaggerInit(lc fx.Lifecycle, routerGroup *server.Group) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			routerGroup.Router.GET("/swagger/*filepath", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				httpSwagger.WrapHandler(w, r)
			})
			return nil
		},
	})
}
