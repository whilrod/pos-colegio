package handler

import (
	"backend-pos-colegio/internal/response"
	"backend-pos-colegio/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type VariantHandler struct {
	service *service.VariantService
}

func NewVariantHandler(service *service.VariantService) *VariantHandler {
	return &VariantHandler{
		service: service,
	}
}

func (h *VariantHandler) List(w http.ResponseWriter, r *http.Request) {
	variants, err := h.service.ListVariants(r.Context())
	if err != nil {
		log.Printf("ERROR GET /variantes: %v", err)
		response.JSONError(w, http.StatusInternalServerError, "Error interno del servidor")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(variants)
}
