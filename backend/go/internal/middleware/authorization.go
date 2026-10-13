package middleware

import (
	"backend-pos-colegio/internal/response"
	"net/http"
)

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, ok := GetAuthenticatedSession(r.Context())
			if !ok {
				response.JSONError(w, http.StatusUnauthorized, "Sesión no autenticada")
				return
			}
			for _, assignedPermision := range session.Permisos {
				if assignedPermision == permission {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.JSONError(w, http.StatusForbidden, "Permiso insuficiente")
		})
	}
}
