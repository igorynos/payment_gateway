package webhookhandler

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/webhook", func(router chi.Router) {
		router.Post("/{provider}", h.UpdateStatusByID)
	})
}
