package models

import "time"

type User struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre" binding:"required,min=3,max=50"`
	Apellido string `json:"apellido" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=6"` // Solo para recibir, no devolver
	Rol      string `json:"rol" binding:"required,oneof=admin doctor recepcion"`

	//Campo especifico para personal medico
	Especialidad string `json:"especialidad,omitempty"` // solo si es medico
	Matricula    string `json:"matricula,omitempty"`    // solo si es medico

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
