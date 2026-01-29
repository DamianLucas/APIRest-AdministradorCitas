package services

import (
	"adminApp/repository"
	"adminApp/utils"
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
		return nil, utils.NewInternal("error al verificar email")
	}
	if usuario == nil {
		return nil, utils.NewUnauthorized("credenciales invalidas")
	}
	if !utils.VerificarPassword(password, usuario.Password) {
		return nil, utils.NewUnauthorized("credenciales inválidas")
	}

	token, err := utils.GenerarToken(usuario.ID, usuario.Rol)
	if err != nil {
		return nil, utils.NewInternal("error al generar token")
	}

	var resp LoginResponse
	resp.Token = token
	resp.User.ID = usuario.ID
	resp.User.Email = usuario.Email
	resp.User.Rol = usuario.Rol

	return &resp, nil

}
