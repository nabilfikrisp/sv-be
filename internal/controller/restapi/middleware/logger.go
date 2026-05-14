// Package middleware provides HTTP middleware.
package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nabilfikrisp/sv-be/pkg/logger"
)

func buildRequestMessage(c *gin.Context) string {
	var result strings.Builder

	result.WriteString(c.ClientIP())
	result.WriteString(" - ")
	result.WriteString(c.Request.Method)
	result.WriteString(" ")
	result.WriteString(c.Request.RequestURI)
	result.WriteString(" - ")
	result.WriteString(strconv.Itoa(c.Writer.Status()))
	result.WriteString(" ")
	result.WriteString(strconv.Itoa(c.Writer.Size()))

	return result.String()
}

// Logger returns a Gin middleware that logs request details after handling.
func Logger(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		l.Info("%s", buildRequestMessage(c))
	}
}
