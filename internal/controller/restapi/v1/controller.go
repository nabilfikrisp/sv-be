package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/nabilfikrisp/sv-be/internal/usecase"
	"github.com/nabilfikrisp/sv-be/pkg/logger"
)

// V1 contains dependencies for version 1 REST API handlers.
type V1 struct {
	uc_post usecase.Post
	l       logger.Interface
	v       *validator.Validate
}
