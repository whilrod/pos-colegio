package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

var (
	ErrProductCodeConflict     = errors.New("el código del producto ya existe")
	ErrProductCategoryNotFound = errors.New("la categoría del producto no existe")
	ErrProductNotFound         = errors.New("el producto no existe")
)

type Product struct {
	IDProducto  string `json:"idproducto"`
	Codigo      string `json:"codigo"`
	Descripcion string `json:"descripcion"`
	IDCategoria string `json:"idcategoria"`
	Categoria   string `json:"categoria"`
	Estado      string `json:"estado"`
}

func (r *Repository) ListProducts(ctx context.Context) ([]Product, error) {
	const query = `
	SELECT
		p.idproducto,
		p.codigo,
		p.descripcion,
		p.idcategoria,
		c.descripcion AS categoria,
		p.estado
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
			&product.Estado,
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
					idcategoria,
					estado
			)
			SELECT
				np.idproducto,
				np.codigo,
				np.descripcion,
				np.idcategoria,
				c.descripcion AS categoria,
				np.estado
			FROM nuevo_producto np
			INNER JOIN categoria c
				ON c.idcategoria = np.idcategoria
		`
	var product Product

	err := r.db.QueryRowContext(ctx, query, codigo, descripcion, idCategoria).Scan(
		&product.IDProducto, &product.Codigo, &product.Descripcion, &product.IDCategoria, &product.Categoria, &product.Estado,
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

func (r *Repository) GetProductByID(ctx context.Context, idProducto string) (*Product, error) {
	const query = `
			SELECT
				p.idproducto,
				p.codigo,
				p.descripcion,
				p.idcategoria,
				c.descripcion AS categoria,
				p.estado
			FROM producto p
			INNER JOIN categoria c
				ON c.idcategoria = p.idcategoria
			WHERE p.idproducto = $1
		`
	var product Product

	err := r.db.QueryRowContext(ctx, query, idProducto).Scan(
		&product.IDProducto,
		&product.Codigo,
		&product.Descripcion,
		&product.IDCategoria,
		&product.Categoria,
		&product.Estado,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProductNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("no se pudo consultar el producto: %w", err)
	}
	return &product, nil
}

func (r *Repository) UpdateProduct(
	ctx context.Context,
	idProducto string,
	codigo string,
	descripcion string,
	idCategoria string,
) (*Product, error) {
	const query = `
		WITH producto_actualizado AS (
			UPDATE producto
			SET
				codigo = $2,
				descripcion = $3,
				idcategoria = $4
			WHERE idproducto = $1
			RETURNING
				idproducto,
				codigo,
				descripcion,
				idcategoria,
				estado
		)
		SELECT
			pa.idproducto,
			pa.codigo,
			pa.descripcion,
			pa.idcategoria,
			c.descripcion AS categoria,
			pa.estado
		FROM producto_actualizado pa
		INNER JOIN categoria c
			ON c.idcategoria = pa.idcategoria
	`

	var product Product

	err := r.db.QueryRowContext(
		ctx,
		query,
		idProducto,
		codigo,
		descripcion,
		idCategoria,
	).Scan(
		&product.IDProducto,
		&product.Codigo,
		&product.Descripcion,
		&product.IDCategoria,
		&product.Categoria,
		&product.Estado,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProductNotFound
	}

	if err != nil {
		var postgresError *pq.Error

		if errors.As(err, &postgresError) {
			switch {
			case postgresError.Code == "23505" &&
				postgresError.Constraint == "uq_producto_codigo":
				return nil, ErrProductCodeConflict

			case postgresError.Code == "23503" &&
				postgresError.Constraint == "fk_producto_categoria":
				return nil, ErrProductCategoryNotFound
			}
		}

		return nil, fmt.Errorf("no se pudo actualizar el producto: %w", err)
	}

	return &product, nil
}

func (r *Repository) DeactivateProduct(
	ctx context.Context,
	idProducto string,
) (*Product, error) {
	const query = `
		WITH producto_desactivado AS (
			UPDATE producto
			SET estado = 'INACTIVO'
			WHERE idproducto = $1
			  AND estado = 'ACTIVO'
			RETURNING
				idproducto,
				codigo,
				descripcion,
				idcategoria,
				estado
		)
		SELECT
			pd.idproducto,
			pd.codigo,
			pd.descripcion,
			pd.idcategoria,
			c.descripcion AS categoria,
			pd.estado
		FROM producto_desactivado pd
		INNER JOIN categoria c
			ON c.idcategoria = pd.idcategoria
	`

	var product Product

	err := r.db.QueryRowContext(ctx, query, idProducto).Scan(
		&product.IDProducto,
		&product.Codigo,
		&product.Descripcion,
		&product.IDCategoria,
		&product.Categoria,
		&product.Estado,
	)

	if errors.Is(err, sql.ErrNoRows) {
		// Diferenciar entre producto inexistente y ya inactivo.
		existingProduct, lookupErr := r.GetProductByID(ctx, idProducto)
		if lookupErr != nil {
			return nil, lookupErr
		}

		if existingProduct.Estado == "INACTIVO" {
			return existingProduct, nil
		}

		return nil, ErrProductNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("no se pudo desactivar el producto: %w", err)
	}

	return &product, nil
}
