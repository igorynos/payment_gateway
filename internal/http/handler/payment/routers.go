package paymenthandler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/payment", func(router chi.Router) {
		router.Post("/", h.Create)
		router.Get("/{paymentID}", h.GetByID)
	})
}
