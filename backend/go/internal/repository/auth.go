package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

type ActiveSession struct {
	IDEmpleado uuid.UUID
	Nombre     string
	Correo     string
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

		return nil, fmt.Errorf(
			"no se pudo consultar el empleado: %w",
			err,
		)
	}

	return &employee, nil
}

func (r *Repository) GetPermissionsByRole(
	ctx context.Context,
	idRol uuid.UUID,
) ([]string, error) {
	const query = `
		SELECT
			p.descripcion
		FROM rol_permiso rp
		JOIN permiso p
			ON p.idpermiso = rp.idpermiso
		WHERE rp.idrol = $1
		ORDER BY p.descripcion
	`

	rows, err := r.db.QueryContext(ctx, query, idRol)
	if err != nil {
		return nil, fmt.Errorf(
			"no se pudieron consultar los permisos del rol: %w",
			err,
		)
	}
	defer rows.Close()

	var permissions []string

	for rows.Next() {
		var permission string

		if err := rows.Scan(&permission); err != nil {
			return nil, fmt.Errorf(
				"no se pudo leer un permiso del rol: %w",
				err,
			)
		}

		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"error al recorrer los permisos del rol: %w",
			err,
		)
	}

	return permissions, nil
}

func (r *Repository) CreateSession(
	ctx context.Context,
	idEmpleado uuid.UUID,
	tokenHash string,
	fechaExpiracion time.Time,
) error {
	const query = `
		INSERT INTO sesion (
			idempleado,
			token_hash,
			fecha_expiracion
		)
		VALUES ($1, $2, $3)
	`

	if _, err := r.db.ExecContext(
		ctx,
		query,
		idEmpleado,
		tokenHash,
		fechaExpiracion,
	); err != nil {
		return fmt.Errorf("no se pudo crear la sesión: %w", err)
	}

	return nil
}

func (r *Repository) GetActiveSessionByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*ActiveSession, error) {
	const query = `
		SELECT
			e.idempleado,
			e.nombre,
			e.correo,
			e.estado,
			e.idrol
		FROM sesion s
		JOIN empleado e
			ON e.idempleado = s.idempleado
		WHERE s.token_hash = $1
			AND s.fecha_revocacion IS NULL
			AND s.fecha_expiracion > CURRENT_TIMESTAMP
	`

	var session ActiveSession

	err := r.db.QueryRowContext(
		ctx,
		query,
		tokenHash,
	).Scan(
		&session.IDEmpleado,
		&session.Nombre,
		&session.Correo,
		&session.Estado,
		&session.IDRol,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}

		return nil, fmt.Errorf(
			"no se pudo validar la sesión: %w",
			err,
		)
	}

	return &session, nil
}
