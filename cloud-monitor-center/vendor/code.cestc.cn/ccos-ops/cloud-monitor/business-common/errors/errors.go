package errors

type BusinessError struct {
	Code    int
	Message string
}

func (e BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(msg string) error {
	return &BusinessError{Code: 500, Message: msg}
}
func NewBusinessErrorCode(code int, msg string) error {
	return &BusinessError{Code: code, Message: msg}
}
