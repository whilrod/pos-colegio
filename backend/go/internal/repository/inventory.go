package repository

import (
	"context"
	"fmt"
)

type Inventory struct {
	IDInventario       string `json:"idinventario"`
	IDVarianteProducto string `json:"idvariante_producto"`
	Codigo             string `json:"codigo"`
	Producto           string `json:"producto"`
	Talla              string `json:"talla"`
	Color              string `json:"color"`
	Genero             string `json:"genero"`
	IDUbicacion        string `json:"idubicacion"`
	Ubicacion          string `json:"ubicacion"`
	Stock              int    `json:"stock"`
}

func (r *Repository) ListInventory(ctx context.Context) ([]Inventory, error) {
	const query = `
			SELECT
				i.idinventario,
				vp.idvariante_producto,
				vp.codigo,
				p.descripcion AS producto,
				t.descripcion AS talla,
				c.descripcion AS color,
				vp.genero,
				u.idubicacion,
				u.descripcion AS ubicacion,
				i.stock
			FROM inventario AS i
			INNER JOIN variante_producto AS vp
				ON vp.idvariante_producto = i.idvariante_producto
			INNER JOIN producto AS p
				ON p.idproducto = vp.idproducto
			INNER JOIN talla AS t
				ON t.idtalla = vp.idtalla
			INNER JOIN color AS c
				ON c.idcolor = vp.idcolor
			INNER JOIN ubicacion AS u
				ON u.idubicacion = i.idubicacion
			ORDER BY p.descripcion, t.descripcion, c.descripcion, u.descripcion
		`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar inventario: %w", err)
	}
	defer rows.Close()

	var inventory []Inventory

	for rows.Next() {
		var item Inventory

		if err := rows.Scan(
			&item.IDInventario,
			&item.IDVarianteProducto,
			&item.Codigo,
			&item.Producto,
			&item.Talla,
			&item.Color,
			&item.Genero,
			&item.IDUbicacion,
			&item.Ubicacion,
			&item.Stock,
		); err != nil {
			return nil, fmt.Errorf("error al leer inventario: %w", err)
		}
		inventory = append(inventory, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al recorrer inventario: %w", err)
	}

	if inventory == nil {
		inventory = []Inventory{}
	}
	return inventory, nil
}
