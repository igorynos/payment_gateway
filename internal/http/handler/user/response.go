package userhandler

import "payment_gateway/internal/user"

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func newUserResponse(value user.User) UserResponse {
	return UserResponse{
		ID:       value.ID,
		Username: value.Username,
		Email:    value.Email,
		Role:     string(value.Role),
	}
}
