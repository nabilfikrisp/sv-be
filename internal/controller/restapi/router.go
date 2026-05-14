// Package restapi configures HTTP routing and REST API handlers.
package restapi

import (
	"github.com/gin-gonic/gin"
	"github.com/nabilfikrisp/sv-be/config"
	"github.com/nabilfikrisp/sv-be/internal/controller/restapi/middleware"
	"github.com/nabilfikrisp/sv-be/internal/usecase"
	"github.com/nabilfikrisp/sv-be/pkg/logger"

	_ "github.com/nabilfikrisp/sv-be/docs" // Swagger Docs
	v1 "github.com/nabilfikrisp/sv-be/internal/controller/restapi/v1"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
//
//	@title       Go Post Article
//	@description API for Go Post Article
//	@version     1.0
//	@host        localhost:8080
//	@BasePath    /v1
func NewRouter(engine *gin.Engine, cfg *config.Config, uc_post usecase.Post, l logger.Interface) {
	// Options
	engine.Use(middleware.Logger(l))
	engine.Use(middleware.Recovery(l))

	// Swagger
	if cfg.Swagger.Enabled {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Routers
	apiV1Group := engine.Group("/v1")
	{
		v1.NewRoutes(*apiV1Group, uc_post, l)
	}
}
