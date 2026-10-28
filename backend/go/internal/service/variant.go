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

var ErrInvalidVariant = errors.New("los datos de la variante no son validos")

type VariantService struct {
	repository *repository.Repository
}

type CreateVariantInput struct {
	IDTalla     string  `json:"idtalla"`
	IDColor     string  `json:"idcolor"`
	Codigo      string  `json:"codigo"`
	Descripcion string  `json:"descripcion"`
	PrecioVenta float64 `json:"precio_venta"`
	Genero      string  `json:"genero"`
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

func (s *VariantService) CreateVariant(ctx context.Context, idProducto string, input CreateVariantInput) (*repository.Variant, error) {
	idProducto = strings.TrimSpace(idProducto)
	input.IDTalla = strings.TrimSpace(input.IDTalla)
	input.IDColor = strings.TrimSpace(input.IDColor)
	input.Codigo = strings.TrimSpace(input.Codigo)
	input.Descripcion = strings.TrimSpace(input.Descripcion)
	input.Genero = strings.ToUpper(strings.TrimSpace(input.Genero))

	productID, err := uuid.Parse(idProducto)
	if err != nil || productID == uuid.Nil {
		return nil, ErrInvalidVariant
	}

	sizeID, err := uuid.Parse(input.IDTalla)
	if err != nil || sizeID == uuid.Nil {
		return nil, ErrInvalidVariant
	}

	colorID, err := uuid.Parse(input.IDColor)
	if err != nil || colorID == uuid.Nil {
		return nil, ErrInvalidVariant
	}

	if input.Codigo == "" ||
		utf8.RuneCountInString(input.Codigo) > 50 ||
		input.Descripcion == "" ||
		utf8.RuneCountInString(input.Descripcion) > 200 ||
		input.PrecioVenta < 0 {
		return nil, ErrInvalidVariant
	}

	variant, err := s.repository.CreateVariant(
		ctx,
		productID.String(),
		sizeID.String(),
		colorID.String(),
		input.Codigo,
		input.Descripcion,
		input.PrecioVenta,
		input.Genero,
	)

	if err != nil {
		return nil, fmt.Errorf("no se pudo crear la variante: %w", err)
	}
	return variant, nil
}
