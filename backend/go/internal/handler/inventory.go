package handler

import (
	"backend-pos-colegio/internal/response"
	"backend-pos-colegio/internal/service"
	"encoding/json"
	"log"
	"net/http"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) List(w http.ResponseWriter, r *http.Request) {
	inventory, err := h.service.List(r.Context())
	if err != nil {
		log.Printf("ERROR GET /inventario: %v", err)
		response.JSONError(w, http.StatusInternalServerError, "Error interno del servidor")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(inventory)
}
