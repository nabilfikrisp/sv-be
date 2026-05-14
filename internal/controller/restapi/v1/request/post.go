// Package request provides HTTP request types.
package request

import (
	"github.com/nabilfikrisp/sv-be/internal/dto"
	"github.com/nabilfikrisp/sv-be/internal/entity"
)

// CreatePost -.
type CreatePost struct {
	Title    string           `json:"title"    validate:"required,min=20,max=200"                           example:"how to learn go programming"`
	Content  string           `json:"content"  validate:"required,min=200"                                   example:"Lorem ipsum dolor sit amet..."`
	Category string           `json:"category" validate:"required,min=3,max=100"                              example:"programming"`
	Status   entity.PostStatus `json:"status"  validate:"required,oneof=publish draft thrash" example:"publish"`
} // @name v1.CreatePost

// UpdatePost -.
type UpdatePost struct {
	Title    *string          `json:"title"    validate:"omitempty,min=20,max=200"                          example:"updated title"`
	Content  *string          `json:"content"  validate:"omitempty,min=200"                                example:"updated content..."`
	Category *string          `json:"category" validate:"omitempty,min=3,max=100"                             example:"updated category"`
	Status   *entity.PostStatus `json:"status"  validate:"omitempty,oneof=publish draft thrash" example:"draft"`
} // @name v1.UpdatePost

// PostFilter -.
type PostFilter struct {
	Status *entity.PostStatus `json:"status"  validate:"omitempty,oneof=publish draft thrash"`
	Limit  *uint64            `json:"limit"`
	Offset *uint64            `json:"offset"`
} // @name v1.PostFilter

// ToDTO converts request filter to DTO.
func (f *PostFilter) ToDTO() dto.PostFilter {
	return dto.PostFilter{
		Status: f.Status,
		Limit:  f.Limit,
		Offset: f.Offset,
	}
}