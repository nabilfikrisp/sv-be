// Package entity defines article domain errors.
package entity

import "errors"

// Post errors.
var (
	ErrPostNotFound        = errors.New("post not found")
	ErrPostStatusInvalid   = errors.New("invalid post status")
	ErrPostTitleEmpty      = errors.New("post title cannot be empty")
	ErrPostAlreadyExists   = errors.New("post already exists")
	ErrInternalServerError = errors.New("internal server error")
)
