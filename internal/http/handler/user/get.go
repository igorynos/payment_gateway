package userhandler

import (
	"errors"
	"net/http"
	"strconv"

	"payment_gateway/internal/lib/api/response"
	"payment_gateway/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type GetByIDResponse struct {
	response.Response
	User UserResponse `json:"user"`
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil || id <= 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid user ID"))
		return
	}

	foundUser, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, user.ErrNotFound) {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, response.Error("user not found"))
		return
	}
	if errors.Is(err, user.ErrInvalidInput) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid user ID"))
		return
	}
	if err != nil {
		h.log.Error("failed to get user", "error", err, "user_id", id)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get user"))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, GetByIDResponse{
		Response: response.OK(),
		User: UserResponse{
			ID:       foundUser.ID,
			Username: foundUser.Username,
			Email:    foundUser.Email,
			Role:     string(foundUser.Role),
		},
	})
}
