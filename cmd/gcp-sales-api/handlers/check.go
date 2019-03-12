package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/ardanlabs/service/internal/platform/web"

	"go.opencensus.io/trace"
)

// health validates the service is healthy and ready to accept requests.
func health(ctx context.Context, log *log.Logger, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "sales-api.handlers.Check.Health")
	defer span.End()

	status := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	web.Respond(ctx, log, w, status, http.StatusOK)

	return nil
}
