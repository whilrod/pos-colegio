package service

import (
	"backend-pos-colegio/internal/repository"
	"context"
	"fmt"
)

type VariantService struct {
	repository *repository.Repository
}

func NewVariantService(repository *repository.Repository) *VariantService {
	return &VariantService{
		repository: repository,
	}
}

func (s *VariantService) ListVariants(ctx context.Context) ([]repository.Variant, error) {
	variants, err := s.repository.ListVariants(ctx)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el listado de variantes: %w", err)
	}
	return variants, nil
}
