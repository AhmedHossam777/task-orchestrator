package middleware

import (
	"log/slog"
	"time"
	
	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Before c.Next(): capture start state ---
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		
		// --- Execute the rest of the chain ---
		c.Next()
		
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		
		requestId := c.GetHeader(RequestIDKey)
		
		slog.Info(
			"request completed", "request_id", requestId,
			"method", method,
			"path", path,
			"status", statusCode,
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
