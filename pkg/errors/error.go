package errors

type ServiceError struct {
	Code    int
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	return e.Message
}

func NewServiceError(code int, message string, err error) *ServiceError {
	return &ServiceError{Code: code, Message: message, Err: err}
}
