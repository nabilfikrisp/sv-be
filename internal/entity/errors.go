// Package entity defines article domain errors.
package entity

import "errors"

// Post errors.
var (
	ErrArticleNotFound      = errors.New("article not found")
	ErrArticleStatusInvalid = errors.New("invalid article status")
	ErrArticleTitleEmpty    = errors.New("article title cannot be empty")
	ErrInternalServerError  = errors.New("internal server error")
)
