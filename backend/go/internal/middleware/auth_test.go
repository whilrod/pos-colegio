package middleware

import (
	"backend-pos-colegio/internal/service"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockSessionValidator struct {
	session *service.AuthenticatedSession
	err     error
}

func (m mockSessionValidator) ValidateSession(
	ctx context.Context,
	token string,
) (*service.AuthenticatedSession, error) {
	return m.session, m.err
}

func TestAuth(t *testing.T) {
	tests := []struct {
		name               string
		authorization      string
		session            *service.AuthenticatedSession
		validationError    error
		expectedStatusCode int
	}{
		{
			name:               "sesión válida",
			authorization:      "Bearer token-valido",
			session:            &service.AuthenticatedSession{},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name:               "token ausente",
			authorization:      "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "esquema de autorización incorrecto",
			authorization:      "Basic token-valido",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "token vacío",
			authorization:      "Bearer ",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "token compuesto solo por espacios",
			authorization:      "Bearer    ",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "sesión inválida",
			authorization:      "Bearer token-invalido",
			validationError:    service.ErrInvalidSession,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "empleado inactivo",
			authorization:      "Bearer token-valido",
			validationError:    service.ErrEmployeeInactive,
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name:               "error inesperado al validar sesión",
			authorization:      "Bearer token-valido",
			validationError:    errors.New("error de base de datos"),
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := mockSessionValidator{
				session: tt.session,
				err:     tt.validationError,
			}

			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})

			handler := Auth(validator)(next)

			request := httptest.NewRequest(
				http.MethodGet,
				"/test",
				nil,
			)
			if tt.authorization != "" {
				request.Header.Set("Authorization", tt.authorization)
			}

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.expectedStatusCode {
				t.Errorf(
					"código HTTP incorrecto: obtenido %d, esperado %d",
					recorder.Code,
					tt.expectedStatusCode,
				)
			}
		})
	}
}
