package handler

import (
	"net/http"

	"backend-pos-colegio/internal/response"
	"backend-pos-colegio/internal/service"
)

type HealthHandler struct {
	service *service.HealthService
}

func NewHealthHandler(service *service.HealthService) *HealthHandler {
	return &HealthHandler{
		service: service,
	}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.Check()
	if err != nil {
		response.JSONError(
			w,
			http.StatusInternalServerError,
			"Error de conexión con la base de datos",
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}
