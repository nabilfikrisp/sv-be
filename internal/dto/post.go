// Package dto provides data transfer objects.
package dto

import "github.com/nabilfikrisp/sv-be/internal/entity"

type (
	// PostFilter defines the structure for filtering posts.
	PostFilter struct {
		Status *entity.PostStatus
		Limit  *uint64
		Offset *uint64
	}

	// PostCreate represents a post creation request.
	PostCreate struct {
		Title    string
		Content  string
		Category string
		Status   entity.PostStatus
	}

	// PostUpdate represents a post update request (Patch style).
	PostUpdate struct {
		Title    *string
		Content  *string
		Category *string
		Status   *entity.PostStatus
	}
)
