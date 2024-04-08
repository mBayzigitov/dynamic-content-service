package auth

import (
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
)

const (
	UnauthorizedUserPrefix  = "uup"
	AuthorizedUserPrefix    = "aup"
	UnauthorizedAdminPrefix = "uap"
	AuthorizedAdminPrefix   = "aap"
)

// ValidateToken
/*
Imitates the authentication service

Returns true if token belongs to administrator, otherwise returns false
*/
func ValidateToken(token string) (bool, error) {

	prefix := token[:3]
	switch prefix {

	case UnauthorizedUserPrefix, UnauthorizedAdminPrefix:
		return false, &serverr.Error{Description: "Пользователь не авторизован", ErrType: serverr.AuthError, HttpStatus: 401}

	case AuthorizedUserPrefix:
		return false, nil

	case AuthorizedAdminPrefix:
		return true, nil

	default:
		return false, &serverr.Error{Description: "Пользователь не авторизован", ErrType: serverr.AuthError, HttpStatus: 401}

	}

}
