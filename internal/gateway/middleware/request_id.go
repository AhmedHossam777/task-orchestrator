package middleware

import (
	"crypto/rand"
	"fmt"
	
	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.GetHeader("X-Request-ID")
		if requestId == "" {
			requestId = generateRequestID()
		}
		c.Set(RequestIDKey, requestId)
		c.Header("X-Request-ID", requestId)
		
		c.Next()
	}
}

func generateRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("req_%x", b)
}
