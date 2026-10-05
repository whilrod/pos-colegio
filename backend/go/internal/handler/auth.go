package handler

import (
	"encoding/json"
	"errors"
	"net/http"

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
