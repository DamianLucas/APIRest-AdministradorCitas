package repository

import (
	"adminApp/internal/database"
	"adminApp/internal/models"
	"database/sql"
)

// Crear Usuario
func CrearUsuario(user *models.User) error {
	query := `
		INSERT INTO users (nombre, apellido, email, password, rol, especialidad, matricula)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at 
	`

	err := database.DB.QueryRow(
		query,
		user.Nombre,
		user.Apellido,
		user.Email,
		user.Password,
		user.Rol,
		user.Especialidad,
		user.Matricula,
	).Scan(&user.ID, &user.CreatedAt)

	return err
}

// ObtenerUsuarioPorEmail - busca un usuario por email
func ObtenerUsuarioPorEmail(email string) (*models.User, error) {
	query := `SELECT id, nombre, apellido, email, password, rol, especialidad, matricula, created_at, updated_at FROM users
		WHERE email = $1 
	`

	var user models.User
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Nombre,
		&user.Apellido,
		&user.Email,
		&user.Password,
		&user.Rol,
		&user.Especialidad,
		&user.Matricula,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get - Obtener usuarios
func ObtenerUsuarios() ([]models.User, error) {
	query := `SELECT id, nombre, apellido, email, rol, created_at, updated_at FROM users ORDER BY id `

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Nombre,
			&user.Apellido,
			&user.Email,
			&user.Rol,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil

}

// Obtener Usuarios por ID
func ObtenerUsuarioPorID(id int) (*models.User, error) {
	query := `SELECT id, nombre, apellido, email, rol, created_at, updated_at FROM users WHERE id = $1`

	var user models.User

	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Nombre,
		&user.Apellido,
		&user.Email,
		&user.Rol,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, err
}

// Actualizar Usuario /:id
func ActualizarUsuario(id int, user *models.User) error {
	query := `
		UPDATE users
		SET nombre = $1, apellido = $2, email = $3, password = $4, updated_at = NOW() 
		WHERE id = $5
	`

	resultado, err := database.DB.Exec(
		query,
		user.Nombre,
		user.Apellido,
		user.Email,
		user.Password,
		id,
	)

	if err != nil {
		return err
	}
	rowsAffected, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Eliminar Usuario
func EliminarUsuario(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	resultado, err := database.DB.Exec(query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := resultado.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
