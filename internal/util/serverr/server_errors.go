package serverr

// error types
const (
	AuthErrorType = "AUTH_ERROR"
)

// defined errors
var (
	UserUnauthorizedError = &ApiError{
		Description: "Пользователь не авторизован",
		ErrType:     AuthErrorType,
		HttpStatus:  401,
	}
	ForbiddenAccessError = &ApiError{
		Description: "Пользователь не имеет доступа",
		ErrType:     AuthErrorType,
		HttpStatus:  403,
	}
)

type ApiError struct {
	Description string `json:"-"`
	ErrType     string `json:"type"`
	HttpStatus  int    `json:"-"`
}

func (e *ApiError) Error() string {
	return e.ErrType + ": " + e.Description
}
