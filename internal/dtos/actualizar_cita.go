package dtos

type ActualizarCitaRequest struct {
	Fecha  string `json:"fecha"`
	Hora   string `json:"hora"`
	Motivo string `json:"motivo"`
	Estado string `json:"estado"`
}
