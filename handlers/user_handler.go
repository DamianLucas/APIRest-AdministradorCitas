package handlers

import (
	"adminApp/models"
	"adminApp/repository"
	"adminApp/services"
	"adminApp/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Crear Usuario
func CrearUsuario(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		utils.BadRequest(c, "Datos invalidos")
		return
	}

	err := services.CrearUsuario(&newUser)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, http.StatusCreated, nil, "usuario creado")
}

// Get - Obtener usuarios
func ObtenerUsuarios(c *gin.Context) {
	users, err := repository.ObtenerUsuarios()
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}
	utils.Success(c, http.StatusOK, gin.H{
		"total":    len(users),
		"usuarios": users,
	}, "")
}

// Obtener Usuarios por ID
func ObtenerUsuarioPorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "ID inválido")
		return
	}

	user, err := services.ObtenerUsuarioPorID(id)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, user, "")
}

// Actualizar Usuario /:id
func ActualizarUsuario(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "ID invalido")
		return
	}

	var datos models.User
	if err := c.ShouldBindJSON(&datos); err != nil {
		utils.BadRequest(c, "datos inválidos")
		return
	}
	err = services.ActualizarUsuario(id, &datos)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, nil, "usuario actualizado")
}

// Eliminar Usuario
func EliminarUsuario(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "ID invalido")
		return
	}

	err = services.EliminarUsuario(id)
	if err != nil {
		utils.HandleServiceError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, nil, "usuario eliminado")
}
