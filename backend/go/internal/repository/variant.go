package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrVariantCodeConflict        = errors.New("el codigo de la variante ya existe")
	ErrVariantCombinationConflict = errors.New("la caombinación de variante ya existe")
	ErrVariantProdctNotFound      = errors.New("el producto no existe")
	ErrVariantSizeNotFound        = errors.New("la talla no existe")
	ErrVariantColorNotFound       = errors.New("el color no existe")
)

type Variant struct {
	IDVarianteProducto string  `json:"idvariante_producto"`
	IDProducto         string  `json:"idproducto"`
	Producto           string  `json:"producto"`
	IDTalla            string  `json:"idtalla"`
	Talla              string  `json:"talla"`
	IDColor            string  `json:"idcolor"`
	Color              string  `json:"color"`
	Codigo             string  `json:"codigo"`
	Descripcion        string  `json:"descripcion"`
	PrecioVenta        float64 `json:"precio_venta"`
	Genero             string  `json:"genero"`
}

func (r *Repository) ListVariants(ctx context.Context) ([]Variant, error) {
	const query = `SELECT
		vp.idvariante_producto,
		vp.idproducto,
		p.descripcion AS producto,
		vp.idtalla,
		t.descripcion AS talla,
		vp.idcolor,
		c.descripcion AS color,
		vp.codigo,
		vp.descripcion,
		vp.precio_venta,
		vp.genero
		FROM variante_producto as vp
		INNER JOIN producto as p
		ON p.idproducto = vp.idproducto
		INNER JOIN talla AS t
		ON t.idtalla = vp.idtalla
		INNER JOIN color AS c
		ON c.idcolor = vp.idcolor
		ORDER BY p.descripcion, t.descripcion, c.descripcion, vp.codigo
		`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar variantes: %w", err)
	}

	defer rows.Close()

	var variants []Variant
	for rows.Next() {
		var variant Variant
		if err := rows.Scan(
			&variant.IDVarianteProducto,
			&variant.IDProducto,
			&variant.Producto,
			&variant.IDTalla,
			&variant.Talla,
			&variant.IDColor,
			&variant.Color,
			&variant.Codigo,
			&variant.Descripcion,
			&variant.PrecioVenta,
			&variant.Genero); err != nil {
			return nil, fmt.Errorf("error al leer variante: %w", err)
		}
		variants = append(variants, variant)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al recorrer variantes: %w", err)
	}

	if variants == nil {
		variants = []Variant{}
	}
	return variants, nil
}

func (r *Repository) CreateVariant(ctx context.Context, idProducto string, idTalla string, idColor string, codigo string, descripcion string, precioVenta float64, genero string) (*Variant, error) {
	const query = `
			WITH nueva_variante AS (
				INSERT INTO variante_producto (
					idproducto,
					idtalla,
					idcolor,
					codigo,
					descripcion,
					precio_venta,
					genero
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING
					idvariante_producto,
					idproducto,
					idtalla,
					idcolor,
					codigo,
					descripcion,
					precio_venta,
					genero
			)
			SELECT
				vp.idvariante_producto,
				vp.idproducto,
				p.descripcion AS producto,
				vp.idtalla,
				t.descripcion AS talla,
				vp.idcolor,
				c.descripcion AS color,
				vp.codigo,
				vp.descripcion,
				vp.precio_venta,
				vp.genero
			FROM nueva_variante AS vp
			INNER JOIN producto AS p
				ON p.idproducto = vp.idproducto
			INNER JOIN talla AS t
				ON t.idtalla = vp.idtalla
			INNER JOIN color AS c
				ON c.idcolor = vp.idcolor
		`
	var variant Variant
	err := r.db.QueryRowContext(
		ctx,
		query,
		idProducto,
		idTalla,
		idColor,
		codigo,
		descripcion,
		precioVenta,
		genero,
	).Scan(
		&variant.IDVarianteProducto,
		&variant.IDProducto,
		&variant.Producto,
		&variant.IDTalla,
		&variant.Talla,
		&variant.IDColor,
		&variant.Color,
		&variant.Codigo,
		&variant.Descripcion,
		&variant.PrecioVenta,
		&variant.Genero,
	)
	if err != nil {
		var postgresError *pq.Error

		if errors.As(err, &postgresError) {
			switch {
			case postgresError.Code == "23505" && postgresError.Constraint == "uq_variante_producto_codigo":
				return nil, ErrVariantCodeConflict
			case postgresError.Code == "23505" && postgresError.Constraint == "uq_variante_producto_combinacion":
				return nil, ErrVariantCombinationConflict
			case postgresError.Code == "23503" && postgresError.Constraint == "fk_variante_producto_producto":
				return nil, ErrVariantProdctNotFound
			case postgresError.Code == "23503" && postgresError.Constraint == "fk_variante_producto_talla":
				return nil, ErrVariantSizeNotFound
			case postgresError.Code == "23503" && postgresError.Constraint == "fk_variante_producto_color":
				return nil, ErrVariantColorNotFound
			}
		}
		return nil, fmt.Errorf("no se pudo crear la variante: %w", err)
	}
	return &variant, nil
}
