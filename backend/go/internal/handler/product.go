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

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idProducto := chi.URLParam(r, "idproducto")

	product, err := h.service.GetProductByID(r.Context(), idProducto)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProduct):
			response.JSONError(w, http.StatusBadRequest, "El identificador del producto no es válido")
		case errors.Is(err, repository.ErrProductNotFound):
			response.JSONError(w, http.StatusNotFound, "El producto indicado no existe")
		default:
			log.Printf("ERROR GET /productos/%s: %v", idProducto, err)
			response.JSONError(w, http.StatusInternalServerError, "Error interno del servidor")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input service.CreateProductInput
	if err := decoder.Decode(&input); err != nil {
		response.JSONError(w, http.StatusBadRequest, "El JSON del producto no es válido")
		return
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		response.JSONError(w, http.StatusBadRequest, "La solicitud debe contener un solo objeto JSON")
		return
	}

	product, err := h.service.CreateProduct(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProduct):
			response.JSONError(w, http.StatusBadRequest, "Código, descripción o categoría inválidos")
		case errors.Is(err, repository.ErrProductCodeConflict):
			response.JSONError(w, http.StatusConflict, "Ya existe un producto con ese código")
		case errors.Is(err, repository.ErrProductCategoryNotFound):
			response.JSONError(w, http.StatusNotFound, "La categoría indicada no existe")
		default:
			log.Printf("ERROR POST /productos: %v", err)
			response.JSONError(w, http.StatusInternalServerError, "Error interno del servidor")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/productos/"+product.IDProducto)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()

	var input service.UpdateProductInput
	if err := decode.Decode(&input); err != nil {
		response.JSONError(w, http.StatusBadRequest, "El JSON del producto no es válido")
		return
	}

	var extra any
	if err := decode.Decode(&extra); err != io.EOF {
		response.JSONError(w, http.StatusBadRequest, "La solicitud debe contener un solo objeto JSON")
		return
	}

	idProducto := chi.URLParam(r, "idproducto")

	product, err := h.service.UpdateProduct(r.Context(), idProducto, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProduct):
			response.JSONError(w, http.StatusBadRequest, "Código, descripción, categoría o identificador inválido")
		case errors.Is(err, repository.ErrProductNotFound):
			response.JSONError(w, http.StatusNotFound, "El producto indicado no existe")
		case errors.Is(err, repository.ErrProductCodeConflict):
			response.JSONError(w, http.StatusConflict, "Ya existe un producto con este código")
		case errors.Is(err, repository.ErrProductCategoryNotFound):
			response.JSONError(w, http.StatusNotFound, "La categoría indicada no existe")
		default:
			log.Printf("ERROR PUT /productos/%s: %v", idProducto, err)
			response.JSONError(w, http.StatusInternalServerError, "Error interno del servidor")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	idProducto := chi.URLParam(r, "idproducto")

	product, err := h.service.DeactivateProduct(r.Context(), idProducto)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidProduct):
			response.JSONError(
				w,
				http.StatusBadRequest,
				"El identificador del producto no es válido",
			)
		case errors.Is(err, repository.ErrProductNotFound):
			response.JSONError(
				w,
				http.StatusNotFound,
				"El producto indicado no existe",
			)
		default:
			log.Printf("ERROR POST /productos/%s/desactivar: %v", idProducto, err)
			response.JSONError(w, http.StatusInternalServerError, "Error interno del servidor")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(product)
}
