package dtos

type CitaResponse struct {
	ID     int    `json:"id"`
	Fecha  string `json:"fecha"`
	Hora   string `json:"hora"`
	Motivo string `json:"motivo,omitempty"`
	Estado string `json:"estado"`

	Paciente PacienteResumen `json:"paciente"`
	Doctor   DoctorResumen   `json:"doctor"`
}

type PacienteResumen struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

type DoctorResumen struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}
