package pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/zelalem-12/distributed-file-storage/internal/server"
	"go.uber.org/fx"
)

func StartHTTPServer(lifecycle fx.Lifecycle, cfg Config, apiRouter *server.Group) {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.SERVER_PORT),
		Handler: apiRouter.Router,
	}

	// Use fx.Lifecycle to manage the server startup and shutdown.
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Starting the server asynchronously
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Could not start server: %s\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping HTTP server...")

			// Gracefully shut down the server with a timeout context.
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				return err
			}

			log.Println("Server stopped cleanly")
			return nil
		},
	})

}
