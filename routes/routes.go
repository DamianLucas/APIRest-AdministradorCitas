package routes

import (
	"adminApp/handlers"
	"adminApp/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// ===== RUTAS PÚBLICAS =====
	r.POST("/login", handlers.Login)

	// ===== RUTAS PROTEGIDAS =====
	api := r.Group("/api")
	api.Use(middleware.RequiereAuth())

	// ===== USUARIOS (solo ADMIN) =====
	usuarios := api.Group("/usuarios")
	usuarios.Use(middleware.RequiereRol("admin"))
	{
		usuarios.POST("/", handlers.CrearUsuario)
		usuarios.GET("/", handlers.ObtenerUsuarios)
		usuarios.GET("/:id", handlers.ObtenerUsuarioPorID)
		usuarios.PUT("/:id", handlers.ActualizarUsuario)
		usuarios.DELETE("/:id", handlers.EliminarUsuario)
	}

	// ===== PACIENTES (ADMIN + RECEPCIÓN + DOCTOR) =====
	pacientes := api.Group("/pacientes")
	pacientes.Use(middleware.RequiereRol("admin", "recepcion", "doctor"))
	{
		pacientes.POST("/", handlers.CrearPaciente)
		pacientes.GET("/", handlers.ListarPacientes)
		pacientes.GET("/:id", handlers.ObtenerPacienteID)
		pacientes.PUT("/:id", handlers.ActualizarPaciente)
		pacientes.DELETE("/:id", handlers.EliminarPaciente)
	}

	// ===== CITAS =====
	citas := api.Group("/citas")
	citas.Use(middleware.RequiereRol("admin", "recepcion", "doctor"))
	{
		citas.POST("/", handlers.CrearCita)
		citas.GET("/", handlers.ObtenerCitas)
		citas.GET("/:id", handlers.ObtenerCitaID)
		citas.PUT("/:id", handlers.ActualizarCita)
		citas.DELETE("/:id", handlers.EliminarCita)
	}
}
