package dtos

type CrearCitaRequest struct {
	PacienteID int    `json:"paciente_id" binding:"required"`
	DoctorID   int    `json:"doctor_id"`
	Fecha      string `json:"fecha" binding:"required"`
	Hora       string `json:"hora" binding:"required"`
	Motivo     string `json:"motivo"`
}
