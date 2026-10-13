package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"backend-pos-colegio/internal/middleware"
	"backend-pos-colegio/internal/response"
	"backend-pos-colegio/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

type LoginRequest struct {
	Correo   string `json:"correo"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.JSONError(
			w,
			http.StatusBadRequest,
			"Solicitud inválida",
		)
		return
	}

	result, err := h.service.Login(
		r.Context(),
		request.Correo,
		request.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.JSONError(
				w,
				http.StatusUnauthorized,
				"Credenciales inválidas",
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(result)
}

func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	session, ok := middleware.GetAuthenticatedSession(r.Context())
	if !ok {
		response.JSONError(w, http.StatusUnauthorized, "Sesion no autenticada")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(session)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const bearerPrefix = "Bearer "

	authorization := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, bearerPrefix) {
		response.JSONError(w, http.StatusUnauthorized, "Token de sesión requerido")
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorization, bearerPrefix))
	if token == "" {
		response.JSONError(w, http.StatusUnauthorized, "Token de sesión requerido")
		return
	}
	if err := h.service.Logout(r.Context(), token); err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidSession):
			response.JSONError(w, http.StatusUnauthorized, "Sesión inválida o expirada")
		default:
			response.JSONError(w, http.StatusInternalServerError, "Error interno del servidor")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
