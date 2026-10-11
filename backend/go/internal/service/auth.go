package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend-pos-colegio/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrEmployeeInactive   = errors.New("empleado inactivo")
	ErrInvalidSession     = errors.New("sesión inválida o expirada")
)

type AuthService struct {
	repository             *repository.Repository
	sessionExpirationHours int
}

func NewAuthService(
	repository *repository.Repository,
	sessionExpirationHours int,
) *AuthService {
	return &AuthService{
		repository:             repository,
		sessionExpirationHours: sessionExpirationHours,
	}
}

type LoginResult struct {
	IDEmpleado string   `json:"idempleado"`
	Nombre     string   `json:"nombre"`
	Correo     string   `json:"correo"`
	IDRol      string   `json:"idrol"`
	Permisos   []string `json:"permisos"`
	Token      string   `json:"token"`
}

type AuthenticatedSession struct {
	IDEmpleado string   `json:"idempleado"`
	Nombre     string   `json:"nombre"`
	Correo     string   `json:"correo"`
	IDRol      string   `json:"idrol"`
	Permisos   []string `json:"permisos"`
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

		return nil, fmt.Errorf(
			"no se pudo autenticar el empleado: %w",
			err,
		)
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

	permissions, err := s.repository.GetPermissionsByRole(
		ctx,
		employee.IDRol,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"no se pudieron obtener los permisos del empleado: %w",
			err,
		)
	}

	token, tokenHash, err := generateSessionToken()
	if err != nil {
		return nil, err
	}

	fechaExpiracion := time.Now().UTC().Add(
		time.Duration(s.sessionExpirationHours) * time.Hour,
	)

	if err := s.repository.CreateSession(
		ctx,
		employee.IDEmpleado,
		tokenHash,
		fechaExpiracion,
	); err != nil {
		return nil, fmt.Errorf(
			"no se pudo crear la sesión del empleado: %w",
			err,
		)
	}

	return &LoginResult{
		IDEmpleado: employee.IDEmpleado.String(),
		Nombre:     employee.Nombre,
		Correo:     employee.Correo,
		IDRol:      employee.IDRol.String(),
		Permisos:   permissions,
		Token:      token,
	}, nil
}

func (s *AuthService) ValidateSession(
	ctx context.Context,
	token string,
) (*AuthenticatedSession, error) {
	token = strings.TrimSpace(token)

	if token == "" {
		return nil, ErrInvalidSession
	}

	tokenHash := hashSessionToken(token)

	session, err := s.repository.GetActiveSessionByTokenHash(
		ctx,
		tokenHash,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidSession
		}

		return nil, fmt.Errorf(
			"no se pudo validar la sesión del empleado: %w",
			err,
		)
	}

	if session.Estado != "ACTIVO" {
		return nil, ErrEmployeeInactive
	}

	permissions, err := s.repository.GetPermissionsByRole(
		ctx,
		session.IDRol,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"no se pudieron consultar los permisos de la sesión: %w",
			err,
		)
	}

	return &AuthenticatedSession{
		IDEmpleado: session.IDEmpleado.String(),
		Nombre:     session.Nombre,
		Correo:     session.Correo,
		IDRol:      session.IDRol.String(),
		Permisos:   permissions,
	}, nil
}

func generateSessionToken() (string, string, error) {
	const tokenBytes = 32

	rawToken := make([]byte, tokenBytes)

	if _, err := rand.Read(rawToken); err != nil {
		return "", "", fmt.Errorf(
			"no se pudo generar el token aleatorio: %w",
			err,
		)
	}

	token := hex.EncodeToString(rawToken)

	return token, hashSessionToken(token), nil
}

func hashSessionToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Evita que un UUID inválido recibido internamente se acepte por accidente
// en futuras operaciones que necesiten identificar un rol.
func parseRoleID(idRol string) (uuid.UUID, error) {
	parsed, err := uuid.Parse(idRol)
	if err != nil {
		return uuid.Nil, fmt.Errorf("idrol inválido: %w", err)
	}

	return parsed, nil
}
