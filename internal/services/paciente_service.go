package services

import (
	"adminApp/internal/models"
	"adminApp/internal/repository"
	"adminApp/pkg/apperrors"
)

func CrearPaciente(p *models.Paciente, rol string, userID int) error {

	// Si es doctor → se asigna a sí mismo
	if rol == "doctor" {
		p.DoctorID = userID
	}

	// Si es recepcion → debe venir doctor_id
	if rol == "recepcion" && p.DoctorID == 0 {
		return apperrors.NewBadRequest("recepcion debe asignar un doctor al paciente")
	}

	err := repository.CrearPaciente(p)
	if err != nil {
		return apperrors.NewInternal("error al crear paciente")
	}

	return nil
}

func ListarPacientes(rol string, userID int) ([]models.Paciente, error) {
	if rol == "doctor" {
		pacientes, err := repository.ListarPacientesPorDoctor(userID)
		if err != nil {
			return nil, apperrors.NewInternal("error al listar pacientes")
		}
		return pacientes, nil
	}

	pacientes, err := repository.ListarPacientes()
	if err != nil {
		return nil, apperrors.NewInternal("error al listar pacientes")
	}
	return pacientes, nil
}

func ObtenerPacienteID(id int, rol string, userID int) (*models.Paciente, error) {
	//validar que cada doctor obtenga su paciente
	if rol == "doctor" {
		ok, err := repository.PacientePerteneceADoctor(id, userID)
		if err != nil {
			return nil, apperrors.NewInternal("error al validar pertenencia")
		}
		if !ok {
			return nil, apperrors.NewForbidden("no puede acceder a pacientes que no son suyos")
		}
	}

	paciente, err := repository.ObtenerPacienteID(id)
	if err != nil {
		return nil, apperrors.NewInternal("error al obtener paciente")
	}
	if paciente == nil {
		return nil, apperrors.NewNotFound("paciente no encontrado")
	}
	return paciente, nil
}

func ActualizarPaciente(id int, data *models.Paciente, rol string, userID int) error {
	//validar si es doctor
	if rol == "doctor" {
		ok, err := repository.PacientePerteneceADoctor(id, userID)
		if err != nil {
			return apperrors.NewInternal("error al validar pertenencia")
		}
		if !ok {
			return apperrors.NewForbidden("no puede modificar pacientes que no son suyos")
		}
	}

	actualizado, err := repository.ActualizarPaciente(id, data)
	if err != nil {
		return apperrors.NewInternal("error al actualizar paciente")
	}
	if !actualizado {
		return apperrors.NewNotFound("paciente no encontrado")
	}

	return nil
}

func EliminarPaciente(id int, rol string, userID int) error {
	//validar si es doctor
	if rol == "doctor" {
		ok, err := repository.PacientePerteneceADoctor(id, userID)
		if err != nil {
			return apperrors.NewInternal("error al validar pertenencia")
		}
		if !ok {
			return apperrors.NewForbidden("no puede eliminar pacientes que no son suyos")
		}
	}

	eliminado, err := repository.EliminarPaciente(id)
	if err != nil {
		return apperrors.NewInternal("error al eliminar paciente")
	}
	if !eliminado {
		return apperrors.NewNotFound("paciente no encontrado")
	}

	return nil
}
