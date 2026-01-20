package handlers

import (
	"adminApp/repository"
	"adminApp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// login
func Login(c *gin.Context) {

	var loginReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos Invalidos",
			"detalle": err.Error(),
		})
		return
	}

	// Buscar usuario por email
	usuario, err := repository.ObtenerUsuarioPorEmail(loginReq.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al verificar email",
			"detalle": err.Error(),
		})
		return
	}

	if usuario == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Credenciales inválidas",
		})
		return
	}

	//verificar password
	if !utils.VerificarPassword(loginReq.Password, usuario.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Credenciales inválidas",
		})
		return
	}

	//generar JWT
	token, err := utils.GenerarToken(usuario.ID, usuario.Rol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al generar token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    usuario.ID,
			"email": usuario.Email,
			"rol":   usuario.Rol,
		},
	})
}
