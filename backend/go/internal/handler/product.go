package handler

import (
	"backend-pos-colegio/internal/response"
	"backend-pos-colegio/internal/service"
	"encoding/json"
	"net/http"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		response.JSONError(w, http.StatusInternalServerError, "Error interno del servidor")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(products)
}
