package service

import "backend-pos-colegio/internal/repository"

type HealthService struct {
	repository *repository.Repository
}

func NewHealthService(repository *repository.Repository) *HealthService {
	return &HealthService{
		repository: repository,
	}
}

func (s *HealthService) Check() (string, error) {
	if err := s.repository.Ping(); err != nil {
		return "", err
	}

	return "OK", nil
}
