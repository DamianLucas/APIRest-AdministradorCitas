package handlers

import (
	"adminApp/internal/dtos"
	"adminApp/internal/services"
	"adminApp/pkg/apperrors"
	"adminApp/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST/Crear Cita
func CrearCita(c *gin.Context) {
	var req dtos.CrearCitaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Datos invalidos")
		return
	}

	//Verifiacar rol
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	err := services.CrearCita(&req, rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, nil, "Cita creada correctamente")
}

// GET/Mostrar Citas
func ObtenerCitas(c *gin.Context) {
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	citas, err := services.ObtenerCitas(rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"total": len(citas),
		"citas": citas,
	}, "")
}

// Cita por id
func ObtenerCitaPorID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID invalido")
		return
	}

	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	cita, err := services.ObtenerCitaPorID(id, rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}
	response.Success(c, http.StatusOK, cita, "")
}

// actualizar cita
func ActualizarCita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID invalido")
		return
	}

	var req dtos.ActualizarCitaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Datos invalidos")
		return
	}

	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	err = services.ActualizarCita(id, &req, rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "Cita actualizada correctamente")
}

// borrar cita
func EliminarCita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Si es doctor, validar que la cita sea suya
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	err = services.EliminarCita(id, rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "Cita eliminada correctamente")
}
