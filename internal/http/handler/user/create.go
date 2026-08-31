package userhandler

import (
	"errors"
	"net/http"
	"payment_gateway/internal/lib/api/response"
	"payment_gateway/internal/user"

	"github.com/go-chi/render"

	"github.com/go-playground/validator/v10"
)

type CreateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type CreateResponse struct {
	response.Response
	User UserResponse `json:"user"`
}

func (h *Handler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("validation failed"))
		return
	}
	createdUser, err := h.service.CreateUser(
		r.Context(),
		user.CreateInput{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if errors.Is(err, user.ErrAlreadyExists) {
		render.Status(r, http.StatusConflict)
		render.JSON(w, r, response.Error("user already exists"))
		return
	}
	if errors.Is(err, user.ErrInvalidInput) {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid user data"))
		return
	}
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to create user"))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateResponse{
		Response: response.OK(),
		User:     newUserResponse(createdUser),
	})
}
