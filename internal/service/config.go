package service

import (
	"context"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/auth"
	"github.com/mBayzigitov/dynamic-content-service/internal/util/serverr"
	"net/http"
)

func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get token from header
		token := r.Header.Get("X-Access-Token")
		if token == "" {
			http.Error(w, serverr.UserUnauthorizedError.JsonBody(), serverr.UserUnauthorizedError.HttpStatus)
			return
		}

		// validate token
		isAdmin, err := auth.ValidateToken(token)
		if err != nil {
			http.Error(w, err.JsonBody(), err.HttpStatus)
			return
		}

		// set isAdmin flag in request context
		ctx := context.WithValue(r.Context(), "isAdmin", isAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
