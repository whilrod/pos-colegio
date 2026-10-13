package middleware

import (
	"backend-pos-colegio/internal/service"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequirePermission(t *testing.T) {
	test := []struct {
		name               string
		session            *service.AuthenticatedSession
		permission         string
		expectedStatusCode int
	}{
		{
			name: "Permiso concedido",
			session: &service.AuthenticatedSession{
				Permisos: []string{"CONSULTAR_PRODUCTO"},
			},
			permission:         "CONSULTAR_PRODUCTO",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "permmiso denegado",
			session: &service.AuthenticatedSession{
				Permisos: []string{"CONSULTAR_CLIENTE"},
			},
			permission:         "CONSULTAR_PRODUCTO",
			expectedStatusCode: http.StatusForbidden,
		},
		{
			name:               "sesión ausente",
			session:            nil,
			permission:         "CONSULTAR_PRODUCTO",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodGet,
				"/test",
				nil,
			)
			if tt.session != nil {
				ctx := context.WithValue(
					request.Context(),
					AuthenticatedSessionKey,
					tt.session,
				)
				request = request.WithContext(ctx)
			}
			recorder := httptest.NewRecorder()
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})

			handler := RequirePermission(tt.permission)(next)
			handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.expectedStatusCode {
				t.Errorf("codigo HTTP incorrecto: obtenido %d, esperado %d", recorder.Code, tt.expectedStatusCode)
			}
		})

	}
}
