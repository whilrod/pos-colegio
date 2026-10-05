package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type EmployeeAuth struct {
	IDEmpleado uuid.UUID
	Nombre     string
	Correo     string
	Password   string
	Estado     string
	IDRol      uuid.UUID
}

func (r *Repository) GetEmployeeByEmail(
	ctx context.Context,
	correo string,
) (*EmployeeAuth, error) {
	const query = `
		SELECT
			idempleado,
			nombre,
			correo,
			password,
			estado,
			idrol
		FROM empleado
		WHERE correo = $1
	`

	var employee EmployeeAuth

	err := r.db.QueryRowContext(ctx, query, correo).Scan(
		&employee.IDEmpleado,
		&employee.Nombre,
		&employee.Correo,
		&employee.Password,
		&employee.Estado,
		&employee.IDRol,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}

		return nil, fmt.Errorf("no se pudo consultar el empleado: %w", err)
	}

	return &employee, nil
}
