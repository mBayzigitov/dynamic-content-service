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
Imitates the authentication handler

Returns true if token belongs to administrator, otherwise returns false
*/
func ValidateToken(token string) (bool, *serverr.ApiError) {

	prefix := token[:3]
	switch prefix {

	case UnauthorizedUserPrefix, UnauthorizedAdminPrefix:
		return false, serverr.UserUnauthorizedError

	case AuthorizedUserPrefix:
		return false, nil

	case AuthorizedAdminPrefix:
		return true, nil

	default:
		return false, serverr.ForbiddenAccessError

	}

}
