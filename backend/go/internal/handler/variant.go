package handler

import (
	"backend-pos-colegio/internal/repository"
	"backend-pos-colegio/internal/response"
	"backend-pos-colegio/internal/service"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (h *VariantHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input service.CreateVariantInput

	if err := decoder.Decode(&input); err != nil {
		response.JSONError(w, http.StatusBadRequest, "El JSON de la variante no es válido")
		return
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		response.JSONError(w, http.StatusBadRequest, "La solicitud debe contener un solo objeto JSON")
		return
	}

	idProducto := chi.URLParam(r, "idproducto")

	variant, err := h.service.CreateVariant(r.Context(), idProducto, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidVariant):
			response.JSONError(w, http.StatusBadRequest, "Los datos de la variante no son válidos")
		case errors.Is(err, repository.ErrVariantCodeConflict):
			response.JSONError(w, http.StatusConflict, "Ya existe una variante con este código")
		case errors.Is(err, repository.ErrVariantProdctNotFound):
			response.JSONError(w, http.StatusNotFound, "El producto indicado no existe")
		case errors.Is(err, repository.ErrVariantSizeNotFound):
			response.JSONError(w, http.StatusNotFound, "La talla indicada no existe")
		case errors.Is(err, repository.ErrVariantColorNotFound):
			response.JSONError(w, http.StatusNotFound, "El color indicado no existe")
		case errors.Is(err, repository.ErrVariantCombinationConflict):
			response.JSONError(w, http.StatusConflict, "Ya existe una variante con esta combinación de producto, talla, color y género")

		default:
			log.Printf("ERROR POST /productos/%s/variantes: %v", idProducto, err)
			response.JSONError(w, http.StatusInternalServerError, "Error interno del servidor")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/variantes/"+variant.IDVarianteProducto)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(variant)
}
