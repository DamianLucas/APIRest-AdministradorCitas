package seed

import (
	"adminApp/internal/models"
	"adminApp/internal/repository"
	"adminApp/pkg/auth"
	"log"
	"os"
)

func SeedAdminUser() {
	// Leer variables de entorno

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	if email == "" || password == "" {
		log.Println("ADMIN_EMAIL o ADMIN_PASSWORD no definidos. Seed cancelado.")
		return
	}

	// Verificar si el admin ya existe
	existing, err := repository.ObtenerUsuarioPorEmail(email)
	if err != nil {
		log.Println("Error al verificart el admin", err)
		return
	}
	if existing != nil {
		log.Println("Admin ya existe")
		return
	}

	//hash password
	hash, err := auth.HashPassword(password)
	if err != nil {
		log.Println("Error al hashear password", err)
		return
	}

	//Crear usuario Admin
	admin := models.User{
		Nombre:   "Administrador",
		Apellido: "Sistema",
		Email:    email,
		Password: hash,
		Rol:      "admin",
	}

	err = repository.CrearUsuario(&admin)
	if err != nil {
		log.Println("Error creando el administrador")
		return
	}

	log.Println("Usuario ADMIN creado correctamente")
}
