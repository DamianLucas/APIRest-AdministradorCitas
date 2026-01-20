package handlers

import (
	"adminApp/models"
	"adminApp/repository"
	"adminApp/utils"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Crear Usuario
func CrearUsuario(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos Invalidos",
			"detalle": err.Error(),
		})
		return
	}
	//validar rol doctor
	if newUser.Rol == "doctor" {
		if newUser.Especialidad == "" || newUser.Matricula == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "doctor debe tener especialidad y matrícula",
			})
			return
		}
	} else {
		newUser.Especialidad = ""
		newUser.Matricula = ""
	}

	// Verificar que el email no exista
	existeEmail, err := repository.ObtenerUsuarioPorEmail(newUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al verificar email",
			"detalle": err.Error(),
		})
		return
	}
	if existeEmail != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "el email ya existe",
		})
		return
	}

	//hashear pass
	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al procesar contraseña",
		})
		return
	}
	newUser.Password = hashedPassword

	//Insertar (LA DB decide si es válido)
	err = repository.CrearUsuario(&newUser)
	if err != nil {

		// interpretar error postgres
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{
					"error": "email o matrícula ya existente",
				})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error al registrar usuario",
		})
		return
	}

	// OK
	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "usuario creado",
	})
}

// Get - Obtener usuarios
func ObtenerUsuarios(c *gin.Context) {
	users, err := repository.ObtenerUsuarios()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al listar usuarios",
			"detalle": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total":    len(users),
		"usuarios": users,
	})
}

// Obtener Usuarios por ID
func ObtenerUsuarioPorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	user, err := repository.ObtenerUsuarioPorID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener Usuario",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Actualizar Usuario /:id
func ActualizarUsuario(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	var datosActualizados models.User
	if err := c.ShouldBindJSON(&datosActualizados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos invalidos",
			"detalle": err.Error(),
		})
		return
	}
	err = repository.ActualizarUsuario(id, &datosActualizados)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al actualizar usuario",
			"detalle": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Usuario actualizado",
	})
}

// Eliminar Usuario
func EliminarUsuario(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	err = repository.EliminarUsuario(id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al eliminar usuario",
			"detalle": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensaje": "Usuario Eliminado"})
}
