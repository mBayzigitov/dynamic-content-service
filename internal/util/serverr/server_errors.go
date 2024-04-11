package serverr

// error types
const (
	UserUnauthorized = "Пользователь не авторизован"
	AccessRestricted = "Пользователь не имеет доступа"
	InvalidData = "Некорректные данные"
)

// defined errors
var (
	UserUnauthorizedError = &ApiError{
		Description: UserUnauthorized,
		HttpStatus:  401,
	}
	ForbiddenAccessError = &ApiError{
		Description: AccessRestricted,
		HttpStatus:  403,
	}
	InvalidRequestError = &ApiError{
		Description: InvalidData,
		ErrType: "Неверный формат запроса",
		HttpStatus: 400,
	}
)

type ApiError struct {
	Description string `json:"description"`
	ErrType     string `json:"error"`
	HttpStatus  int    `json:"-"`
}

func NewInvalidRequestError(errm string) *ApiError {
	return &ApiError{
		Description: "Некорректные данные",
		ErrType: errm,
		HttpStatus: 401,
	}
}

func (e *ApiError) Error() string {
	return e.ErrType + ": " + e.Description
}
