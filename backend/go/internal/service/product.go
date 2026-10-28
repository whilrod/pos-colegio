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

type UpdateProductInput = CreateProductInput

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

func (s *ProductService) GetProductByID(ctx context.Context, idProducto string) (*repository.Product, error) {
	idProducto = strings.TrimSpace(idProducto)

	productID, err := uuid.Parse(idProducto)
	if err != nil || productID == uuid.Nil {
		return nil, ErrInvalidProduct
	}
	product, err := s.repository.GetProductByID(ctx, productID.String())
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el producto: %w", err)
	}

	return product, nil
}

func validateProductInput(input CreateProductInput) (string, string, string, error) {
	codigo := strings.TrimSpace(input.Codigo)
	descripcion := strings.TrimSpace(input.Descripcion)
	idCategoria := strings.TrimSpace(input.IDCategoria)

	if codigo == "" || utf8.RuneCountInString(codigo) > 50 || descripcion == "" || utf8.RuneCountInString(descripcion) > 200 {
		return "", "", "", ErrInvalidProduct
	}
	categoryID, err := uuid.Parse(idCategoria)
	if err != nil || categoryID == uuid.Nil {
		return "", "", "", ErrInvalidProduct
	}
	return codigo, descripcion, categoryID.String(), nil
}

func (s *ProductService) UpdateProduct(
	ctx context.Context,
	idProducto string,
	input UpdateProductInput,
) (*repository.Product, error) {
	idProducto = strings.TrimSpace(idProducto)

	productID, err := uuid.Parse(idProducto)
	if err != nil || productID == uuid.Nil {
		return nil, ErrInvalidProduct
	}

	codigo, descripcion, categoriaID, err := validateProductInput(input)
	if err != nil {
		return nil, err
	}

	product, err := s.repository.UpdateProduct(
		ctx,
		productID.String(),
		codigo,
		descripcion,
		categoriaID,
	)
	if err != nil {
		return nil, fmt.Errorf("no se pudo actualizar el producto: %w", err)
	}

	return product, nil
}
func (s *ProductService) DeactivateProduct(
	ctx context.Context,
	idProducto string,
) (*repository.Product, error) {
	idProducto = strings.TrimSpace(idProducto)

	productID, err := uuid.Parse(idProducto)
	if err != nil || productID == uuid.Nil {
		return nil, ErrInvalidProduct
	}

	product, err := s.repository.DeactivateProduct(
		ctx,
		productID.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("no se pudo desactivar el producto: %w", err)
	}

	return product, nil
}
