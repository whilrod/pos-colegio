package service

import (
	"backend-pos-colegio/internal/repository"
	"context"
)

type InventoryService struct {
	repository *repository.Repository
}

func NewInventoryService(repository *repository.Repository) *InventoryService {
	return &InventoryService{repository: repository}
}

func (s *InventoryService) List(ctx context.Context) ([]repository.Inventory, error) {
	return s.repository.ListInventory(ctx)
}
