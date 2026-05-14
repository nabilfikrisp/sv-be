// Package response provides HTTP response types.
package response

// Error -.
type Error struct {
	Error string `json:"error" example:"message"`
} // @name v1.Error
