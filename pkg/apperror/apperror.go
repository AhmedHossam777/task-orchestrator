package apperror

import "fmt"

type Apperror struct {
	Status  int
	Code    string
	Message string
}

func New(status int, code string, message string) *Apperror {
	return &Apperror{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func (e *Apperror) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Status, e.Code, e.Message)
}

func BadRequest(code string, message string) *Apperror {
	return New(400, code, message)
}

func NotFound(code string, message string) *Apperror {
	return New(404, code, message)
}

func Conflict(code string, message string) *Apperror {
	return New(409, code, message)
}

func Internal(code string, message string) *Apperror {
	return New(500, code, message)
}
