package services

import (
	"adminApp/dtos"
	"adminApp/models"
	"adminApp/repository"
	"adminApp/utils"
	"time"
)

func CrearCita(req *dtos.CrearCitaRequest, rol string, userID int) error {

	doctorID := req.DoctorID

	if rol == "doctor" {
		doctorID = userID
	} else if doctorID == 0 {
		return utils.NewBadRequest("doctor_id es obligatorio")
	}

	fecha, err := time.Parse("2006-01-02", req.Fecha)
	if err != nil {
		return utils.NewBadRequest("formato de fecha inválido (YYYY-MM-DD)")
	}

	if rol == "doctor" {
		ok, err := repository.PacientePerteneceADoctor(req.PacienteID, userID)
		if err != nil {
			return utils.NewInternal("error al validar paciente")
		}
		if !ok {
			return utils.NewForbidden("no puede crear citas para pacientes que no son suyos")
		}
	}

	cita := &models.Cita{
		PacienteID: req.PacienteID,
		DoctorID:   doctorID,
		Fecha:      fecha,
		Hora:       req.Hora,
		Motivo:     req.Motivo,
		Estado:     "pendiente",
	}

	if err := repository.CrearCita(cita); err != nil {
		return utils.NewInternal("error al crear cita")
	}

	return nil

}

func ObtenerCitas(rol string, userID int) ([]dtos.CitaResponse, error) {
	var (
		citas []dtos.CitaResponse
		err   error
	)

	if rol == "doctor" {
		citas, err = repository.ObtenerCitasPorDoctor(userID)
	} else {
		citas, err = repository.ObtenerCitas()
	}

	if err != nil {
		return nil, utils.NewInternal("error al obtener citas")
	}

	return citas, nil
}

func ObtenerCitaPorID(id int, rol string, userID int) (*dtos.CitaResponse, error) {
	// un doctor solo puede ver sus citas
	if rol == "doctor" {
		ok, err := repository.CitaPerteneceADoctor(id, userID)
		if err != nil {
			return nil, utils.NewInternal("error al validar permisos")
		}

		if !ok {
			return nil, utils.NewForbidden("No puede acceder a esta cita")
		}
	}

	cita, err := repository.ObtenerCitaPorID(id)
	if err != nil {
		return nil, utils.NewInternal("error al obtener la cita")
	}
	if cita == nil {
		return nil, utils.NewNotFound("Cita no encontrada")
	}

	return cita, nil
}

func ActualizarCita(id int, req *dtos.ActualizarCitaRequest, rol string, userID int) error {
	//doctor solo puede modificar sus citas
	if rol == "doctor" {
		ok, err := repository.CitaPerteneceADoctor(id, userID)
		if err != nil {
			return utils.NewInternal("error al validar permisos")
		}
		if !ok {
			return utils.NewForbidden("No puede modificar citas que no son suyas")
		}
	}

	//Obtener modelo real desde DB
	citaActual, err := repository.ObtenerCitaModeloPorID(id)
	if err != nil {
		return utils.NewInternal("error al obtener cita")
	}
	if citaActual == nil {
		utils.NewNotFound("Cita no encontrada")
	}

	//Aplicar patch parcial (solo campos enviados)
	if req.Fecha != "" {
		if fecha, err := time.Parse("2006-01-02", req.Fecha); err != nil {
			return utils.NewBadRequest("Formato de fecha inválido (YYYY-MM-DD)")
		} else {
			citaActual.Fecha = fecha
		}
	}
	if req.Hora != "" {
		citaActual.Hora = req.Hora
	}
	if req.Motivo != "" {
		citaActual.Motivo = req.Motivo
	}
	if req.Estado != "" {
		citaActual.Estado = req.Estado
	}

	//guardar cambios
	ok, err := repository.ActualizarCita(id, citaActual)
	if err != nil {
		return utils.NewInternal("error la actualizar cita")
	}
	if !ok {
		return utils.NewNotFound("Cita no encontrada")
	}

	return nil

}

func EliminarCita(id int, rol string, userID int) error {
	if rol == "doctor" {
		ok, err := repository.CitaPerteneceADoctor(id, userID)
		if err != nil {
			return utils.NewInternal("error al validar permisos")
		}
		if !ok {
			return utils.NewForbidden("No puede eliminar citas que no son suyas")
		}
	}

	eliminada, err := repository.EliminarCita(id)
	if err != nil {
		return utils.NewNotFound("error al eliminar cita")
	}

	if !eliminada {
		return utils.NewNotFound("Cita no encontrada")
	}

	return nil

}
