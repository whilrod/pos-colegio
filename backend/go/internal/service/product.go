package service

import (
	"backend-pos-colegio/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrInvalidProduct = errors.New("los datos del producto no son válidos")

type ProductService struct {
	repository *repository.Repository
}

type CreateProductInput struct {
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	IDCategoria string `json:"idcategoria"`
}

func NewProductService(repository *repository.Repository) *ProductService {
	return &ProductService{
		repository: repository,
	}
}

func (s *ProductService) ListProducts(ctx context.Context) ([]repository.Product, error) {
	products, err := s.repository.ListProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el listado de productos: %w", err)
	}
	return products, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, input CreateProductInput) (*repository.Product, error) {
	input.Codigo = strings.TrimSpace(input.Codigo)
	input.Descripcion = strings.TrimSpace(input.Descripcion)
	input.IDCategoria = strings.TrimSpace(input.IDCategoria)

	if input.Codigo == "" || utf8.RuneCountInString(input.Codigo) > 50 || input.Descripcion == "" || utf8.RuneCountInString(input.Descripcion) > 200 {
		return nil, ErrInvalidProduct
	}

	categoryID, err := uuid.Parse(input.IDCategoria)
	if err != nil || categoryID == uuid.Nil {
		return nil, ErrInvalidProduct
	}

	product, err := s.repository.CreateProduct(ctx, input.Codigo, input.Descripcion, categoryID.String())
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear el producto: %w", err)
	}
	return product, nil
}
