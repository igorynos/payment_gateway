package user

import (
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user already exists")
	ErrInvalidInput  = errors.New("invalid user input")
)

type Role string

const (
	RoleUser  Role = "customer"
	RoleAdmin Role = "admin"
)

type User struct {
	ID        int64
	Username  string
	Email     string
	Role      Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateInput struct {
	Username string
	Email    string
	Password string
}

type CreateParams struct {
	Username     string
	Email        string
	PasswordHash string
	Role         Role
}
