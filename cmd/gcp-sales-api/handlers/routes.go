package handlers

import (
	"log"
	"net/http"

	"github.com/ardanlabs/service/internal/mid"
	"github.com/ardanlabs/service/internal/platform/gcp/ds"
	"github.com/ardanlabs/service/internal/platform/web"
)

// API returns a handler for a set of routes.
func API(log *log.Logger, cDatastore *ds.DS, apiEnv string) http.Handler {
	app := web.New(log, mid.ErrorHandler, mid.RequestLogger)

	// Register authentication endpoints.
	switch apiEnv {
	case "development":
		//Register health check endpoint. This route is not authenticated.
		app.Handle("GET", "/dev/sales-api/v1/health", health)

	case "production":
		//Register health check endpoint. This route is not authenticated.
		app.Handle("GET", "/sales-api/v1/health", health)

	}

	return app
}
