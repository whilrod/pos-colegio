package repository

import (
	"context"
	"fmt"
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
