package repository

import "errors"

var (
	ErrNotFound      = errors.New("url not found")
	ErrAleradyExists = errors.New("url already exists")
)
