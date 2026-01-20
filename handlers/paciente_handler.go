package handlers

import (
	"adminApp/models"
	"adminApp/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// POST /pacientes - crear
func CrearPaciente(c *gin.Context) {
	var nuevoPaciente models.Paciente

	if err := c.ShouldBindJSON(&nuevoPaciente); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos Invalidos",
			"detalle": err.Error(),
		})
		return
	}

	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	// Si es doctor → se asigna a sí mismo
	if rol == "doctor" {
		nuevoPaciente.DoctorID = userID
	}

	// Si es recepcion → debe venir doctor_id
	if rol == "recepcion" && nuevoPaciente.DoctorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "recepción debe asignar un doctor al paciente",
		})
		return
	}

	// // (Opcional) validar rol permitido
	// if rol != "doctor" && rol != "recepcion" && rol != "admin" {
	// 	c.JSON(http.StatusForbidden, gin.H{
	// 		"error": "rol no autorizado",
	// 	})
	// 	return
	// }

	err := repository.CrearPaciente(&nuevoPaciente)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al registrar paciente",
			"detalle": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje":  "Paciente registrado",
		"paciente": nuevoPaciente,
	})
}

// GET Pacientes - listar pacientes
func ListarPacientes(c *gin.Context) {
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	var (
		pacientes []models.Paciente
		err       error
	)

	if rol == "doctor" {
		pacientes, err = repository.ListarPacientesPorDoctor(userID)
	} else {
		pacientes, err = repository.ListarPacientes()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al listar pacientes",
			"detalle": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     len(pacientes),
		"pacientes": pacientes,
	})
}

// GET /paciente/:id - Obtener por ID
func ObtenerPacienteID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalido"})
		return
	}

	//validar que cada doctor obtenga su paciente
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	if rol == "doctor" {
		ok, err := repository.PacientePerteneceADoctor(id, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al validar permisos",
			})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No puede acceder a pacientes que no son suyos",
			})
			return
		}
	}

	paciente, err := repository.ObtenerPacienteID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener paciente",
		})
		return
	}
	if paciente == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Paciente no encontrado"})
		return
	}

	c.JSON(http.StatusOK, paciente)
}

// PUT /pacientes/:id - Actualizar
func ActualizarPaciente(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	var datosActualizados models.Paciente
	if err := c.ShouldBindJSON(&datosActualizados); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"detalle": err.Error(),
		})
		return
	}

	//datos del  usuario
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	//validar si es doctor
	if rol == "doctor" {
		ok, err := repository.PacientePerteneceADoctor(id, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al validar permisos",
			})
			return
		} else if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No puede modificar pacientes que no son suyos",
			})
			return
		}
	}

	actualizado, err := repository.ActualizarPaciente(id, &datosActualizados)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al actualizar paciente",
			"detalle": err.Error(),
		})
		return
	}

	if !actualizado {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Paciente no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Paciente actualizado correctamente",
	})
}

// Eliminar Paciente
func EliminarPaciente(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	//validar si es doctor

	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	if rol == "doctor" {
		ok, err := repository.PacientePerteneceADoctor(id, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al validar permisos",
			})
			return
		} else if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No puede eliminar pacientes que no sean suyos",
			})
			return
		}
	}

	eliminado, err := repository.EliminarPaciente(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al eliminar paciente",
			"detalle": err.Error(),
		})
		return
	}

	if !eliminado {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Paciente no encontrado",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Paciente eliminado correctamente",
	})
}
