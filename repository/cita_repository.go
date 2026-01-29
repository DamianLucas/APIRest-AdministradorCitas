package repository

import (
	"adminApp/database"
	"adminApp/dtos"
	"adminApp/models"
	"database/sql"
	"time"
)

// Crear Cita
func CrearCita(cita *models.Cita) error {
	query := `
		INSERT INTO citas
			(paciente_id, doctor_id, fecha, hora, motivo, estado, observaciones)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)	
		RETURNING id
	`
	err := database.DB.QueryRow(
		query,
		cita.PacienteID,
		cita.DoctorID,
		cita.Fecha,
		cita.Hora,
		cita.Motivo,
		cita.Estado,
		cita.Observaciones,
	).Scan(&cita.ID)

	return err
}

// Listar todas (admin/recepción)
func ObtenerCitas() ([]dtos.CitaResponse, error) {
	query := `
		SELECT
			c.id,
			c.fecha,
			c.hora,
			c.motivo,
			c.estado,
			p.id,
			p.nombre,
			u.id,
			u.nombre
		FROM citas c
		JOIN pacientes p ON p.id = c.paciente_id
		JOIN users u ON u.id = c.doctor_id
		ORDER BY c.fecha DESC	
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var citas []dtos.CitaResponse

	for rows.Next() {
		var c dtos.CitaResponse
		var fecha time.Time
		err := rows.Scan(
			&c.ID,
			&fecha,
			&c.Hora,
			&c.Motivo,
			&c.Estado,
			&c.Paciente.ID,
			&c.Paciente.Nombre,
			&c.Doctor.ID,
			&c.Doctor.Nombre,
		)

		if err != nil {
			return nil, err
		}

		c.Fecha = fecha.Format("2006-01-02")
		citas = append(citas, c)
	}

	return citas, err
}

// Para admin/recepción: obtiene cualquier cita
func ObtenerCitaPorID(id int) (*dtos.CitaResponse, error) {
	query := `
	SELECT
		c.id,
		c.fecha,
		c.hora,
		c.motivo,
		c.estado,
		p.id,
		p.nombre,
		p.apellido,
		u.id,
		u.nombre,
		u.apellido
	FROM citas c
	JOIN pacientes p ON p.id = c.paciente_id
	JOIN users u ON u.id = c.doctor_id
	WHERE c.id = $1
	`

	var resp dtos.CitaResponse
	var fecha time.Time

	err := database.DB.QueryRow(query, id).Scan(
		&resp.ID,
		&fecha,
		&resp.Hora,
		&resp.Motivo,
		&resp.Estado,
		&resp.Paciente.ID,
		&resp.Paciente.Nombre,
		&resp.Paciente.Apellido,
		&resp.Doctor.ID,
		&resp.Doctor.Nombre,
		&resp.Doctor.Apellido,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	resp.Fecha = fecha.Format("2006-01-02")
	return &resp, nil
}

func ObtenerCitaModeloPorID(id int) (*models.Cita, error) {
	query := `
		SELECT
			id,
			paciente_id,
			doctor_id,
			fecha,
			hora,
			motivo,
			estado,
			observaciones
		FROM citas
		WHERE id = $1
	`

	var c models.Cita

	err := database.DB.QueryRow(query, id).Scan(
		&c.ID,
		&c.PacienteID,
		&c.DoctorID,
		&c.Fecha,
		&c.Hora,
		&c.Motivo,
		&c.Estado,
		&c.Observaciones,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// Actualizar cita
func ActualizarCita(id int, cita *models.Cita) (bool, error) {
	query := `
		UPDATE citas
		SET
			fecha = $1,
			hora = $2,
			motivo = $3,
			estado = $4,
			updated_at = NOW()
		WHERE id = $5
	`

	result, err := database.DB.Exec(
		query,
		cita.Fecha,
		cita.Hora,
		cita.Motivo,
		cita.Estado,
		id,
	)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

// Eliminar cita
func EliminarCita(id int) (bool, error) {
	query := `DELETE FROM citas WHERE id = $1`

	result, err := database.DB.Exec(query, id)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil // la cita no existe
	}

	return true, nil
}

// ============================================
// FUNCIONES EXTRAS PARA SEGURIDAD
// ============================================

// Doctor: Listar citas de un doctor específico
func ObtenerCitasPorDoctor(doctorID int) ([]dtos.CitaResponse, error) {
	query := `
	SELECT
		c.id,
		c.fecha,
		c.hora,
		c.motivo,
		c.estado,
		p.id,
		p.nombre,
		p.apellido,
		u.id,
		u.nombre,
		u.apellido
		FROM citas c
		JOIN pacientes p ON p.id = c.paciente_id
		JOIN users u ON u.id = c.doctor_id
		WHERE c.doctor_id = $1
		ORDER BY c.fecha DESC, c.hora
	`

	rows, err := database.DB.Query(query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var citas []dtos.CitaResponse

	for rows.Next() {
		var c dtos.CitaResponse

		err := rows.Scan(
			&c.ID,
			&c.Fecha,
			&c.Hora,
			&c.Motivo,
			&c.Estado,
			&c.Paciente.ID,
			&c.Paciente.Nombre,
			&c.Paciente.Apellido,
			&c.Doctor.ID,
			&c.Doctor.Nombre,
			&c.Doctor.Apellido,
		)
		if err != nil {
			return nil, err
		}

		citas = append(citas, c)
	}

	return citas, nil
}

func CitaPerteneceADoctor(citaID int, doctorID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM citas
			WHERE id = $1 AND doctor_id = $2
		)
	`

	var exists bool
	err := database.DB.QueryRow(query, citaID, doctorID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
