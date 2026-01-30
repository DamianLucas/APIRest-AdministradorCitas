package routes

import (
	"adminApp/internal/handlers"
	"adminApp/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// ===== RUTAS PÚBLICAS =====
	r.POST("/login", handlers.Login)

	// ===== RUTAS PROTEGIDAS =====
	api := r.Group("/api")
	api.Use(middleware.RequiereAuth())
	v1 := api.Group("/v1")

	// ===== USUARIOS (solo ADMIN) =====
	usuarios := v1.Group("/usuarios")
	usuarios.Use(middleware.RequiereRol("admin"))
	{
		usuarios.POST("/", handlers.CrearUsuario)
		usuarios.GET("/", handlers.ObtenerUsuarios)
		usuarios.GET("/:id", handlers.ObtenerUsuarioPorID)
		usuarios.PUT("/:id", handlers.ActualizarUsuario)
		usuarios.DELETE("/:id", handlers.EliminarUsuario)
	}

	// ===== PACIENTES (ADMIN + RECEPCIÓN + DOCTOR) =====
	pacientes := v1.Group("/pacientes")
	pacientes.Use(middleware.RequiereRol("admin", "recepcion", "doctor"))
	{
		pacientes.POST("/", handlers.CrearPaciente)
		pacientes.GET("/", handlers.ListarPacientes)
		pacientes.GET("/:id", handlers.ObtenerPacienteID)
		pacientes.PUT("/:id", handlers.ActualizarPaciente)
		pacientes.DELETE("/:id", handlers.EliminarPaciente)
	}

	// ===== CITAS =====
	citas := v1.Group("/citas")
	citas.Use(middleware.RequiereRol("admin", "recepcion", "doctor"))
	{
		citas.POST("/", handlers.CrearCita)
		citas.GET("/", handlers.ObtenerCitas)
		citas.GET("/:id", handlers.ObtenerCitaPorID)
		citas.PUT("/:id", handlers.ActualizarCita)
		citas.DELETE("/:id", handlers.EliminarCita)
	}
}
