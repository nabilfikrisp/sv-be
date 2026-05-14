// Package entity defines article domain models.
package entity

import "time"

// PostStatus represents the status of an article post.
type PostStatus string

// Post status types based on requirements.
const (
	StatusPublish PostStatus = "Publish"
	StatusDraft   PostStatus = "Draft"
	StatusThrash  PostStatus = "Thrash"
)

// Post represents a post entity in the article database.
type Post struct {
	ID          int        `json:"id"           example:"1"`
	Title       string     `json:"title"        example:"How to Learn Go"`
	Content     string     `json:"content"      example:"The content of the article goes here..."`
	Category    string     `json:"category"     example:"Programming"`
	CreatedDate time.Time  `json:"created_date" example:"2026-05-14T11:38:11Z"`
	UpdatedDate time.Time  `json:"updated_date" example:"2026-05-14T11:38:11Z"`
	Status      PostStatus `json:"status"       example:"Publish"`
} // @name entity.Post

// Valid returns true if the post status is valid.
func (s PostStatus) Valid() bool {
	switch s {
	case StatusPublish, StatusDraft, StatusThrash:
		return true
	}
	return false
}
