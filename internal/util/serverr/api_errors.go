package serverr

import "encoding/json"

// error types
const (
	UserUnauthorized = "Пользователь не авторизован"
	AccessRestricted = "Пользователь не имеет доступа"
	InvalidData    = "Некорректные данные"
	ServerConflict = "Внутреннняя ошибка сервера"
	BannerNotFound = "Баннер не найден"
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
		ErrType:     "Неверный формат запроса",
		HttpStatus:  400,
	}
	StorageError = &ApiError{
		Description: ServerConflict,
		ErrType:     "Ошибка хранилища данных",
		HttpStatus:  500,
	}
	TokenParsingError = &ApiError{
		Description: ServerConflict,
		ErrType:     "Ошибка парсинга токена",
		HttpStatus:  500,
	}
	BannerNotFoundError = &ApiError{
		Description: BannerNotFound,
		HttpStatus:  404,
	}
)

type ApiError struct {
	Description string `json:"description"`
	ErrType     string `json:"error"`
	HttpStatus  int    `json:"-"`
}

func (apierr *ApiError) JsonBody() string {
	var res []byte

	if apierr.ErrType != "" {
		res, _ = json.Marshal(struct {
			Description string `json:"description"`
			ErrType     string `json:"error"`
		}{
			Description: apierr.Description,
			ErrType:     apierr.ErrType,
		})
	} else {
		res, _ = json.Marshal(struct {
			Description string `json:"description"`
		}{
			Description: apierr.Description,
		})
	}

	return string(res)
}

func NewInvalidRequestError(errm string) *ApiError {
	return &ApiError{
		Description: "Некорректные данные",
		ErrType:     errm,
		HttpStatus:  400,
	}
}

func (e *ApiError) Error() string {
	return e.ErrType + ": " + e.Description
}
