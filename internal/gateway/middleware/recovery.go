package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	
	"github.com/gin-gonic/gin"
	
	"github.com/AhmedHossam777/task-orchestrator/pkg/response"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get(RequestIDKey)
				
				slog.Error(
					"panic recovered",
					"request_id", requestID,
					"error", err,
					"stack", string(debug.Stack()),
				)
				
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					response.Envelope{
						Success: false,
						Error: &response.Error{
							Code:    "INTERNAL_ERROR",
							Message: "an unexpected error occurred",
						},
					},
				)
			}
		}()
		
		c.Next()
	}
}
