package service

import (
	"backend-pos-colegio/internal/repository"
	"context"
	"fmt"
)

type ProductService struct {
	repository *repository.Repository
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
