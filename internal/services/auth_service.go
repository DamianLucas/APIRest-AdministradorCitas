package services

import (
	"adminApp/internal/repository"
	"adminApp/pkg/apperrors"
	"adminApp/pkg/auth"
)

type LoginResponse struct {
	Token string
	User  struct {
		ID    int
		Email string
		Rol   string
	}
}

func Login(email string, password string) (*LoginResponse, error) {
	usuario, err := repository.ObtenerUsuarioPorEmail(email)
	if err != nil {
		return nil, apperrors.NewInternal("error al verificar email")
	}
	if usuario == nil {
		return nil, apperrors.NewUnauthorized("credenciales invalidas")
	}
	if !auth.VerificarPassword(password, usuario.Password) {
		return nil, apperrors.NewUnauthorized("credenciales inválidas")
	}

	token, err := auth.GenerarToken(usuario.ID, usuario.Rol)
	if err != nil {
		return nil, apperrors.NewInternal("error al generar token")
	}

	var resp LoginResponse
	resp.Token = token
	resp.User.ID = usuario.ID
	resp.User.Email = usuario.Email
	resp.User.Rol = usuario.Rol

	return &resp, nil

}
