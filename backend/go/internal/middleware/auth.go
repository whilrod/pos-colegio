package middleware

import (
	"backend-pos-colegio/internal/response"
	"backend-pos-colegio/internal/service"
	"context"
	"errors"
	"net/http"
	"strings"
)

type AuthContextKey string

const AuthenticatedSessionKey AuthContextKey = "authenticated_session"

type SessionValidator interface {
	ValidateSession(
		ctx context.Context,
		token string,
	) (*service.AuthenticatedSession, error)
}

func Auth(authService SessionValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const bearerPrefix = "Bearer "

			authorization := r.Header.Get("Authorization")
			if !strings.HasPrefix(authorization, bearerPrefix) {
				response.JSONError(
					w,
					http.StatusUnauthorized,
					"Token de sesión requerido",
				)
				return
			}

			token := strings.TrimSpace(
				strings.TrimPrefix(authorization, bearerPrefix),
			)
			if token == "" {
				response.JSONError(w, http.StatusUnauthorized, "Token de sesión requerido")
				return
			}
			session, err := authService.ValidateSession(
				r.Context(),
				token,
			)
			if err != nil {
				switch {
				case errors.Is(err, service.ErrInvalidSession):
					response.JSONError(
						w,
						http.StatusUnauthorized,
						"Sesión inválida o expirada",
					)

				case errors.Is(err, service.ErrEmployeeInactive):
					response.JSONError(
						w,
						http.StatusForbidden,
						"Empleado inactivo",
					)

				default:
					response.JSONError(
						w,
						http.StatusInternalServerError,
						"Error interno del servidor",
					)
				}
				return
			}
			ctx := context.WithValue(
				r.Context(),
				AuthenticatedSessionKey,
				session,
			)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetAuthenticatedSession(
	ctx context.Context,
) (*service.AuthenticatedSession, bool) {
	session, ok := ctx.Value(AuthenticatedSessionKey).(*service.AuthenticatedSession)
	return session, ok && session != nil
}
