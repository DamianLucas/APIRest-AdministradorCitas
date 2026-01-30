package services

import (
	"adminApp/internal/models"
	"adminApp/internal/repository"
	"adminApp/pkg/apperrors"
	"adminApp/pkg/auth"
	"database/sql"

	"github.com/lib/pq"
)

func CrearUsuario(user *models.User) error {
	//validacion de negocio
	if user.Rol == "doctor" {
		if user.Especialidad == "" || user.Matricula == "" {
			return apperrors.NewBadRequest("doctor debe tener especialidad y matrícula")
		}
	} else {
		user.Especialidad = ""
		user.Matricula = ""
	}

	//verificar email
	existeEmail, err := repository.ObtenerUsuarioPorEmail(user.Email)
	if err != nil {
		return apperrors.NewInternal("error al verificar email")
	}
	if existeEmail != nil {
		return apperrors.NewConflict("el email ya existe")
	}

	//hash password
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return apperrors.NewInternal("error al procesar password")
	}
	user.Password = hashedPassword

	//crear usuario
	err = repository.CrearUsuario(user)
	if err == nil {
		return nil
	}

	// Si hay error, verificar si es duplicado
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return apperrors.NewConflict("email o matrícula ya existente")
	}

	return apperrors.NewInternal("error al crear usuario")
}

func ObtenerUsuarios() ([]models.User, error) {
	users, err := repository.ObtenerUsuarios()
	if err != nil {
		return nil, apperrors.NewInternal("error al obtener usuarios")
	}
	return users, nil
}

func ObtenerUsuarioPorID(id int) (*models.User, error) {
	user, err := repository.ObtenerUsuarioPorID(id)
	if err != nil {
		return nil, apperrors.NewInternal("error al obtener usuario")
	}
	if user == nil {
		return nil, apperrors.NewInternal("usuario no encontrado")
	}
	return user, nil
}

func ActualizarUsuario(id int, data *models.User) error {
	err := repository.ActualizarUsuario(id, data)
	if err == sql.ErrNoRows {
		return apperrors.NewNotFound("usuario no encontrado")
	}

	if err != nil {
		return apperrors.NewInternal("error al actualizar usuario")
	}

	return nil
}

func EliminarUsuario(id int) error {
	err := repository.EliminarUsuario(id)

	if err == sql.ErrNoRows {
		return apperrors.NewNotFound("usuario no encontrado")
	}

	if err != nil {
		return apperrors.NewInternal("error al eliminar usuario")
	}

	return nil
}
