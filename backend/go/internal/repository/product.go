package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrProductCodeConflict     = errors.New("el código del producto ya existe")
	ErrProductCategoryNotFound = errors.New("la categoría del producto no existe")
)

type Product struct {
	IDProducto  string `json:"idproducto"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	IDCategoria string `json:"idcategoria"`
	Categoria   string `json:"categoria"`
}

func (r *Repository) ListProducts(ctx context.Context) ([]Product, error) {
	const query = `
	SELECT
		p.idproducto,
		p.codigo,
		p.descripcion,
		p.idcategoria,
		c.descripcion AS categoria
	FROM producto p
	INNER JOIN categoria c
		ON c.idcategoria = p.idcategoria
	ORDER BY p.codigo
`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron consultar los productos: %w", err)
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.IDProducto,
			&product.Codigo,
			&product.Descripcion,
			&product.IDCategoria,
			&product.Categoria,
		); err != nil {
			return nil, fmt.Errorf("no se pudo leer un producto: %w", err)
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al recorrer los productos: %w", err)
	}
	return products, nil
}

func (r *Repository) CreateProduct(ctx context.Context, codigo string, descripcion string, idCategoria string) (*Product, error) {
	const query = `
			WITH nuevo_producto AS (
				INSERT INTO producto (
					codigo,
					descripcion,
					idcategoria
				)
				VALUES ($1, $2, $3)
				RETURNING
					idproducto,
					codigo,
					descripcion,
					idcategoria
			)
			SELECT
				np.idproducto,
				np.codigo,
				np.descripcion,
				np.idcategoria,
				c.descripcion AS categoria
			FROM nuevo_producto np
			INNER JOIN categoria c
				ON c.idcategoria = np.idcategoria
		`
	var product Product

	err := r.db.QueryRowContext(ctx, query, codigo, descripcion, idCategoria).Scan(
		&product.IDProducto, &product.Codigo, &product.Descripcion, &product.IDCategoria, &product.Categoria,
	)
	if err != nil {
		var postgresError *pq.Error

		if errors.As(err, &postgresError) {
			switch {
			case postgresError.Code == "23505" && postgresError.Constraint == "uq_producto_codigo":
				return nil, ErrProductCodeConflict
			case postgresError.Code == "23503" && postgresError.Constraint == "fk_producto_categoria":
				return nil, ErrProductCategoryNotFound
			}
		}
		return nil, fmt.Errorf("no se pudo crear el producto: %w", err)
	}
	return &product, nil
}
