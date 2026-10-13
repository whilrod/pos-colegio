package repository

import (
	"context"
	"fmt"
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
