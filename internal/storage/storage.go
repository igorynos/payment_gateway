package storage

import "errors"

var (
	ErrUserNotFoud = errors.New("User not found")
	ErrUserExist   = errors.New("User already exist")
	ErrEmptyUrl    = errors.New("Empty url")
)
