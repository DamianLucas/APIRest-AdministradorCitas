package services

import (
	"adminApp/models"
	"adminApp/repository"
	"adminApp/utils"
	"database/sql"

	"github.com/lib/pq"
)

func CrearUsuario(user *models.User) error {
	//validacion de negocio
	if user.Rol == "doctor" {
		if user.Especialidad == "" || user.Matricula == "" {
			return utils.NewBadRequest("doctor debe tener especialidad y matrícula")
		}
	} else {
		user.Especialidad = ""
		user.Matricula = ""
	}

	//verificar email
	existeEmail, err := repository.ObtenerUsuarioPorEmail(user.Email)
	if err != nil {
		return utils.NewInternal("error al verificar email")
	}
	if existeEmail != nil {
		return utils.NewConflict("el email ya existe")
	}

	//hash password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return utils.NewInternal("error al procesar password")
	}
	user.Password = hashedPassword

	//crear usuario
	err = repository.CrearUsuario(user)
	if err == nil {
		return nil
	}

	// Si hay error, verificar si es duplicado
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return utils.NewConflict("email o matrícula ya existente")
	}

	return utils.NewInternal("error al crear usuario")
}

func ObtenerUsuarios() ([]models.User, error) {
	users, err := repository.ObtenerUsuarios()
	if err != nil {
		return nil, utils.NewInternal("error al obtener usuarios")
	}
	return users, nil
}

func ObtenerUsuarioPorID(id int) (*models.User, error) {
	user, err := repository.ObtenerUsuarioPorID(id)
	if err != nil {
		return nil, utils.NewInternal("error al obtener usuario")
	}
	if user == nil {
		return nil, utils.NewInternal("usuario no encontrado")
	}
	return user, nil
}

func ActualizarUsuario(id int, data *models.User) error {
	err := repository.ActualizarUsuario(id, data)
	if err == sql.ErrNoRows {
		return utils.NewNotFound("usuario no encontrado")
	}

	if err != nil {
		return utils.NewInternal("error al actualizar usuario")
	}

	return nil
}

func EliminarUsuario(id int) error {
	err := repository.EliminarUsuario(id)

	if err == sql.ErrNoRows {
		return utils.NewNotFound("usuario no encontrado")
	}

	if err != nil {
		return utils.NewInternal("error al eliminar usuario")
	}

	return nil
}
