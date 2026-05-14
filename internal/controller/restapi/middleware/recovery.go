package middleware

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nabilfikrisp/sv-be/pkg/logger"
)

func buildPanicMessage(c *gin.Context, err any) string {
	var result strings.Builder

	result.WriteString(c.ClientIP())
	result.WriteString(" - ")
	result.WriteString(c.Request.Method)
	result.WriteString(" ")
	result.WriteString(c.Request.RequestURI)
	result.WriteString(" PANIC DETECTED: ")
	fmt.Fprintf(&result, "%v\n%s\n", err, debug.Stack())

	return result.String()
}

// Recovery returns a Gin middleware that recovers from panics and logs them.
func Recovery(l logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Error(buildPanicMessage(c, err))
				c.AbortWithStatus(500)
			}
		}()

		c.Next()
	}
}
