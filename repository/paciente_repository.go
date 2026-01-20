package repository

import (
	"adminApp/database"
	"adminApp/models"
	"database/sql"
)

// Crear Y guardar Paciente
func CrearPaciente(paciente *models.Paciente) error {
	query := `
		INSERT INTO pacientes 
			(nombre, apellido, dni, fecha_nacimiento, telefono, email, direccion, obra_social, doctor_id)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at;
	`

	err := database.DB.QueryRow(
		query,
		paciente.Nombre,
		paciente.Apellido,
		paciente.DNI,
		paciente.FechaNacimiento,
		paciente.Telefono,
		paciente.Email,
		paciente.Direccion,
		paciente.ObraSocial,
		paciente.DoctorID,
	).Scan(&paciente.ID, &paciente.CreatedAt, &paciente.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// Mostrar lista de Pacientes
func ListarPacientes() ([]models.Paciente, error) {
	query := `SELECT id, nombre, apellido, dni, fecha_nacimiento, telefono, email, direccion, obra_social, created_at, updated_at FROM pacientes ORDER BY id`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pacientes []models.Paciente
	for rows.Next() {
		var paciente models.Paciente

		err := rows.Scan(
			&paciente.ID,
			&paciente.Nombre,
			&paciente.Apellido,
			&paciente.DNI,
			&paciente.FechaNacimiento,
			&paciente.Telefono,
			&paciente.Email,
			&paciente.Direccion,
			&paciente.ObraSocial,
			&paciente.CreatedAt,
			&paciente.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		pacientes = append(pacientes, paciente)
	}
	return pacientes, nil
}

// Obtener paciente por ID, paciente/id
func ObtenerPacienteID(id int) (*models.Paciente, error) {
	query := `SELECT id, nombre, apellido, dni, fecha_nacimiento, telefono, email, direccion, obra_social, created_at, updated_at FROM pacientes WHERE id = $1`

	var paciente models.Paciente

	err := database.DB.QueryRow(query, id).Scan(
		&paciente.ID,
		&paciente.Nombre,
		&paciente.Apellido,
		&paciente.DNI,
		&paciente.FechaNacimiento,
		&paciente.Telefono,
		&paciente.Email,
		&paciente.Direccion,
		&paciente.ObraSocial,
		&paciente.CreatedAt,
		&paciente.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &paciente, nil
}

// PUT Actualizar paciente por ID
func ActualizarPaciente(id int, paciente *models.Paciente) (bool, error) {
	query := `
		UPDATE pacientes
		SET 
			nombre = $1,
			apellido = $2,
			dni = $3,
			fecha_nacimiento = $4,
			telefono = $5,
			email = $6,
			direccion = $7,
			obra_social = $8,
			updated_at = NOW()
		WHERE id = $9
	`

	result, err := database.DB.Exec(
		query,
		paciente.Nombre,
		paciente.Apellido,
		paciente.DNI,
		paciente.FechaNacimiento,
		paciente.Telefono,
		paciente.Email,
		paciente.Direccion,
		paciente.ObraSocial,
		id,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil // no existe
	}

	return true, nil
}

// DELETE paciente por ID
func EliminarPaciente(id int) (bool, error) {
	query := `DELETE FROM pacientes WHERE id = $1`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil // no existe
	}

	return true, nil
}

// Extra endpoint para doctores
func ListarPacientesPorDoctor(doctorID int) ([]models.Paciente, error) {
	query := `
		SELECT 
			id,
			nombre,
			apellido,
			dni,
			fecha_nacimiento,
			telefono,
			email,
			direccion,
			obra_social,
			doctor_id,
			created_at,
			updated_at
		FROM pacientes
		WHERE doctor_id = $1
		ORDER BY apellido, nombre
	`

	rows, err := database.DB.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// IMPORTANTE: slice inicializado (evita null en JSON)
	pacientes := make([]models.Paciente, 0)

	for rows.Next() {
		var p models.Paciente

		err := rows.Scan(
			&p.ID,
			&p.Nombre,
			&p.Apellido,
			&p.DNI,
			&p.FechaNacimiento,
			&p.Telefono,
			&p.Email,
			&p.Direccion,
			&p.ObraSocial,
			&p.DoctorID,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		pacientes = append(pacientes, p)
	}

	// Verificar errores del cursor
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pacientes, nil
}

func PacientePerteneceADoctor(pacienteID int, doctorID int) (bool, error) {
	query := `
		SELECT EXISTS (
    	SELECT 1
    	FROM pacientes
    	WHERE id = $1
      	AND doctor_id = $2
		);
	`

	var exist bool
	err := database.DB.QueryRow(query, pacienteID, doctorID).Scan(&exist)
	if err != nil {
		return false, nil
	}

	return exist, nil
}
