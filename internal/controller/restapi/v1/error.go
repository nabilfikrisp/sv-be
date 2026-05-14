package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func errorResponse(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"error": msg,
	})
}

// // validationErrorResponse writes a 400 response with field-level validation messages, or a generic bad request error if the error is not a ValidationErrors type.
func validationErrorResponse(c *gin.Context, err error) {
	if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
		messages := make([]string, 0, len(validationErrors))

		for _, e := range validationErrors {
			messages = append(messages, formatFieldError(e))
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": messages,
		})
		return
	}

	errorResponse(c, http.StatusBadRequest, err.Error())
}

// FormatValidationError formats validation errors into user-friendly messages
func FormatValidationError(err error) string {
	var messages []string

	if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
		for _, e := range validationErrors {
			messages = append(messages, formatFieldError(e))
		}
	}

	return strings.Join(messages, "; ")
}

// formatFieldError formats a single validation error into a user-friendly message.
func formatFieldError(e validator.FieldError) string {
	field := e.Field()

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("'%s' is required", field)
	case "oneof":
		return fmt.Sprintf("'%s' must be one of [%s]", field, e.Param())
	case "email":
		return fmt.Sprintf("'%s' must be a valid email", field)
	case "min":
		return fmt.Sprintf("'%s' must be at least %s characters", field, e.Param())
	case "max":
		return fmt.Sprintf("'%s' must be at most %s characters", field, e.Param())
	default:
		return fmt.Sprintf("'%s' failed validation on '%s'", field, e.Tag())
	}
}
