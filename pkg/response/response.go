package response

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(c *gin.Context, data any) {
	c.JSON(
		http.StatusOK, Envelope{
			Success: true,
			Data:    data,
		},
	)
}

func Created(c *gin.Context, data any) {
	c.JSON(
		http.StatusCreated, &Envelope{
			Success: true,
			Data:    data,
		},
	)
}

func Fail(c *gin.Context, httpStatus int, code string, message string) {
	c.JSON(
		httpStatus, &Envelope{
			Success: false,
			Error: &Error{
				Code:    code,
				Message: message,
			},
		},
	)
}
