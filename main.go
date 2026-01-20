package main

import (
	"adminApp/database"
	"adminApp/routes"
	"adminApp/seed"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

/*
 PROBAR CRUD citas
 chequear detelles de endpoints
 PROBAR LAS CORRECIONES DEL CRUD CITAS (PROBABLES FALLOS) <=====


 {
    "cita": {
        "id": 1,
        "paciente_id": 18,
        "paciente": {
            "id": 0,
            "nombre": "",
            "apellido": "",
            "dni": "",
            "fecha_nacimiento": "",
            "telefono": "",
            "email": "",
            "direccion": "",
            "obra_social": "",
            "doctor_id": 0,
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z"
        },
        "doctor_id": 12,
        "doctor": {
            "id": 0,
            "nombre": "",
            "apellido": "",
            "email": "",
            "rol": "",
            "created_at": "0001-01-01T00:00:00Z",
            "updated_at": "0001-01-01T00:00:00Z"
        },
        "fecha": "2026-01-20T00:00:00Z",
        "hora": "15:30",
        "motivo": "Control general",
        "estado": "pendiente",
        "created_at": "0001-01-01T00:00:00Z",
        "updated_at": "0001-01-01T00:00:00Z"
    },
    "mensaje": "Cita creada correctamente"
}

*/

// ===(Opcional más adelante)===

// Refresh tokens
// Logout
// Auditoría (created_by)
// Paginación
// Filtros
// Swagger / documentación

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error al cargar archivo .env")
	}

	database.Connect()
	defer database.Close()

	seed.SeedAdminUser()

	//Rutas
	r := gin.Default()
	routes.SetupRoutes(r)

	// Leer puerto del .env
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}

	// Iniciar servidor
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%s\n", port)
	fmt.Println("🐘 Conectado a PostgreSQL")

	r.Run(":" + port)
}
