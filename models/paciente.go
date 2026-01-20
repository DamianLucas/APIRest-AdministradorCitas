package models

import "time"

type Paciente struct {
	ID              int       `json:"id"`
	Nombre          string    `json:"nombre" binding:"required,min=3,max=200"`
	Apellido        string    `json:"apellido" binding:"required,min=3,max=200"`
	DNI             string    `json:"dni" binding:"required"`
	FechaNacimiento string    `json:"fecha_nacimiento" binding:"required"`
	Telefono        string    `json:"telefono" binding:"required,min=5,max=50"`
	Email           string    `json:"email" binding:"required,email"`
	Direccion       string    `json:"direccion" binding:"required,min=5,max=200"`
	ObraSocial      string    `json:"obra_social" binding:"required,min=3,max=200"`
	DoctorID        int       `json:"doctor_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
