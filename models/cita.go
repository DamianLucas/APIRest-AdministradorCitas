package models

import "time"

type Cita struct {
	ID int `json:"id"`

	//Relaciones
	PacienteID int      `json:"paciente_id" binding:"required"`
	Paciente   Paciente `json:"paciente"`

	DoctorID int  `json:"doctor_id" `
	Doctor   User `json:"doctor"`

	//Datos de la cita
	Fecha         time.Time `json:"fecha" binding:"required"`
	Hora          string    `json:"hora" binding:"required"`
	Motivo        string    `json:"motivo"`
	Estado        string    `json:"estado"`
	Observaciones string    `json:"observaciones,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
