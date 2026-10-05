package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"backend-pos-colegio/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrEmployeeInactive   = errors.New("empleado inactivo")
)

type AuthService struct {
	repository *repository.Repository
}

func NewAuthService(repository *repository.Repository) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

type LoginResult struct {
	IDEmpleado string `json:"idempleado"`
	Nombre     string `json:"nombre"`
	Correo     string `json:"correo"`
	IDRol      string `json:"idrol"`
}

func (s *AuthService) Login(
	ctx context.Context,
	correo string,
	password string,
) (*LoginResult, error) {
	correo = strings.TrimSpace(correo)

	if correo == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	employee, err := s.repository.GetEmployeeByEmail(ctx, correo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf("no se pudo autenticar el empleado: %w", err)
	}

	if employee.Estado != "ACTIVO" {
		return nil, ErrEmployeeInactive
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(employee.Password),
		[]byte(password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &LoginResult{
		IDEmpleado: employee.IDEmpleado.String(),
		Nombre:     employee.Nombre,
		Correo:     employee.Correo,
		IDRol:      employee.IDRol.String(),
	}, nil
}
