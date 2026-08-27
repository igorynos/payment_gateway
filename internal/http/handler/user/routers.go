package userhandler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/users", func(router chi.Router) {
		router.Post("/", h.Create)
		router.Get("/{userID}", h.GetByID)
	})
}
