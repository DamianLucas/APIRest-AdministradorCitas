package handlers

import (
	"adminApp/internal/models"
	"adminApp/internal/services"
	"adminApp/pkg/apperrors"
	"adminApp/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /pacientes - crear
func CrearPaciente(c *gin.Context) {
	var paciente models.Paciente

	if err := c.ShouldBindJSON(&paciente); err != nil {
		response.BadRequest(c, "Datos invalidos")
		return
	}

	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	err := services.CrearPaciente(&paciente, rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, paciente, "Paciente registrado")
}

// GET Pacientes - listar pacientes
func ListarPacientes(c *gin.Context) {
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	pacientes, err := services.ListarPacientes(rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"total":     len(pacientes),
		"pacientes": pacientes,
	}, "")
}

// GET /paciente/:id - Obtener por ID
func ObtenerPacienteID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID invalido")
		return
	}

	//validar que cada doctor obtenga su paciente
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	paciente, err := services.ObtenerPacienteID(id, rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, paciente, "")
}

// PUT /pacientes/:id - Actualizar
func ActualizarPaciente(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID invalido")
		return
	}

	var data models.Paciente
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "Datos invalidos")
		return
	}

	//datos del  usuario
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	err = services.ActualizarPaciente(id, &data, rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "Paciente actualizado correctamente")
}

// Eliminar Paciente
func EliminarPaciente(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	err = services.EliminarPaciente(id, rol, userID)
	if err != nil {
		apperrors.HandleServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, nil, "Paciente eliminado correctamente")
}
