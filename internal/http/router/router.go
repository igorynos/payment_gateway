package router

import (
	"log/slog"
	"net/http"
	mvlogger "payment_gateway/internal/http/middleware/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RouteRegistrar interface {
	RegisterRoutes(router chi.Router)
}

func New(
	log *slog.Logger,
	registrars ...RouteRegistrar,
) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mvlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Route("/api/v1", func(api chi.Router) {
		for _, registrar := range registrars {
			registrar.RegisterRoutes(api)
		}
	})
	return router
}
