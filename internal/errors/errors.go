package errors

type ErrorCode string

const (
	ErrInvalidRequest ErrorCode = "invalid_request"
	ErrInternalServer ErrorCode = "internal_server_error"
	ErrNotFound       ErrorCode = "not_found"
	ErrUnauthorized   ErrorCode = "unauthorized"
	ErrForbidden      ErrorCode = "forbidden"
	ErrConflict       ErrorCode = "conflict"
)

type CustomError struct {
	Code    ErrorCode
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(code ErrorCode, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}
