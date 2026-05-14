package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nabilfikrisp/sv-be/internal/usecase"
	"github.com/nabilfikrisp/sv-be/pkg/logger"
)

// NewRoutes -.
func NewRoutes(apiV1Group gin.RouterGroup, uc_post usecase.Post, l logger.Interface) {
	r := &V1{
		uc_post: uc_post,
		l:       l,
		v:       validator.New(validator.WithRequiredStructEnabled()),
	}

	// Contact
	contactGroup := apiV1Group.Group("/article")
	{
		contactGroup.POST("", r.createPost)
		contactGroup.GET("", r.listPosts)
		contactGroup.GET("/:id", r.getPostByID)
		contactGroup.PATCH("/:id", r.updatePost)
		contactGroup.DELETE("/:id", r.deletePost)
	}
}
