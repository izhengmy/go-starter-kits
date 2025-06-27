package errorx

type ServiceError struct {
	Code    int
	Message string
}

var _ error = (*ServiceError)(nil)

func NewServiceError(message string) *ServiceError {
	return &ServiceError{
		Message: message,
		Code:    0,
	}
}

func (e *ServiceError) Error() string {
	return e.Message
}

func (e *ServiceError) WithCode(code int) *ServiceError {
	e.Code = code
	return e
}
