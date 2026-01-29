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
		port = "SERVER_PORT"
	}

	// Iniciar servidor
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%s\n", port)
	fmt.Println("🐘 Conectado a PostgreSQL")

	r.Run(":" + port)
}
