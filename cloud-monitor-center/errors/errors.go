package errors

import "fmt"

type BusinessError struct {
	Code    int
	Message string
}

func (e BusinessError) Error() string {
	return fmt.Sprintf("code: %q msg:%s", e.Code, e.Message)
}

func NewBussinessError(code int, msg string) error {
	return &BusinessError{Code: code, Message: msg}
}
