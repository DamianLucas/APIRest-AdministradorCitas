package handlers

import (
	"adminApp/dtos"
	"adminApp/models"
	"adminApp/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// POST/Crear Cita
func CrearCita(c *gin.Context) {
	var req dtos.CrearCitaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"detalle": err.Error(),
		})
		return
	}

	//Verifiacar rol
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	doctorID := req.DoctorID
	if rol == "doctor" {
		doctorID = userID
	} else if doctorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "doctor_id es obligatorio",
		})
		return
	}

	// Parsear fecha
	fecha, err := time.Parse("2006-01-02", req.Fecha) //preguntar porque esa fecha en string
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Formato de fecha inválido (YYYY-MM-DD)",
		})
		return
	}

	//modelo real
	nuevaCita := models.Cita{
		PacienteID: req.PacienteID,
		DoctorID:   doctorID,
		Fecha:      fecha,
		Hora:       req.Hora,
		Motivo:     req.Motivo,
		Estado:     "pendiente",
	}

	//Guardar en DB
	if err := repository.CrearCita(&nuevaCita); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al crear la cita",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"mensaje": "Cita creada correctamente",
		"id":      nuevaCita.ID,
	})
}

// GET/Mostrar Citas
func ObtenerCitas(c *gin.Context) {
	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	var (
		citas []dtos.CitaResponse
		err   error
	)

	if rol == "doctor" {
		// El doctor solo ve sus citas
		citas, err = repository.ObtenerCitasPorDoctor(userID)
	} else {
		// Admin y recepción ven todas
		citas, err = repository.ObtenerCitas()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al listar citas",
			"detalle": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(citas),
		"citas": citas,
	})
}

// Cita por id
func ObtenerCitaID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	if rol == "doctor" {
		ok, err := repository.CitaPerteneceADoctor(id, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de permisos"})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No puede acceder a esta cita"})
			return
		}
	}

	cita, err := repository.ObtenerCitaPorID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la cita"})
		return
	}
	if cita == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cita no encontrada"})
		return
	}

	c.JSON(http.StatusOK, cita)
}

// actualizar cita
func ActualizarCita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	var req dtos.ActualizarCitaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos inválidos",
			"detalle": err.Error(),
		})
		return
	}

	rol := c.GetString("rol")
	userID := c.GetInt("user_id")

	// 1️⃣ Validar permisos si es doctor
	if rol == "doctor" {
		ok, err := repository.CitaPerteneceADoctor(id, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al validar permisos",
			})
			return
		}
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No puede modificar citas que no son suyas",
			})
			return
		}
	}

	// 2️⃣ Obtener la cita actual (estado real en DB)
	citaActual, err := repository.ObtenerCitaModeloPorID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener la cita",
		})
		return
	}
	if citaActual == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cita no encontrada",
		})
		return
	}

	// 3️⃣ Aplicar SOLO los campos que vinieron en el request

	if req.Fecha != "" {
		fecha, err := time.Parse("2006-01-02", req.Fecha)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Formato de fecha inválido (YYYY-MM-DD)",
			})
			return
		}
		citaActual.Fecha = fecha
	}

	if req.Hora != "" {
		citaActual.Hora = req.Hora
	}

	if req.Motivo != "" {
		citaActual.Motivo = req.Motivo
	}

	if req.Estado != "" {
		citaActual.Estado = req.Estado
	}

	// 4️⃣ Actualizar en DB
	actualizado, err := repository.ActualizarCita(id, citaActual)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al actualizar la cita",
		})
		return
	}

	if !actualizado {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cita no encontrada",
		})
		return
	}

	// 5️⃣ Respuesta
	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Cita actualizada correctamente",
	})
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
	if rol == "doctor" {
		ok, err := repository.CitaPerteneceADoctor(id, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error al validar permisos",
			})
			return
		}

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "No puede eliminar citas que no son suyas",
			})
			return
		}
	}

	eliminada, err := repository.EliminarCita(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al eliminar cita",
		})
		return
	}

	if !eliminada {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Cita no encontrada",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "Cita eliminada correctamente",
	})
}
