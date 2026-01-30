package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	// leer variables de entorno
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	//Convertir puertro a int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal("❌ DB_PORT debe ser un número")
	}

	//Validar que existan
	if host == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("❌ Faltan variables de entorno de la base de datos")
	}

	//Crear string de conexion
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Reutilizar err (con = no con :=)
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("❌ Error al conectar a la base de datos:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Error al hacer ping a la base de datos:", err)
	}

	fmt.Println("✅ Conectado a PostgreSQL")

	//Crear tabla si no existe
	crearTabla := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		nombre VARCHAR(100) NOT NULL,
		apellido VARCHAR(100) NOT NULL,
		email VARCHAR(200) UNIQUE NOT NULL,
		password VARCHAR (255) NOT NULL,
		rol VARCHAR(50) NOT NULL CHECK (rol IN ('admin', 'doctor', 'recepcion')),
		especialidad VARCHAR(100),
		matricula VARCHAR(50),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS pacientes (
		id SERIAL PRIMARY KEY,
		nombre VARCHAR(100) NOT NULL,
		apellido VARCHAR(100) NOT NULL,
		dni VARCHAR(20) UNIQUE NOT NULL,
		fecha_nacimiento DATE,
		telefono VARCHAR(200),
		email VARCHAR(200),
		direccion TEXT,
		obra_social VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS citas (
		id SERIAL PRIMARY KEY,
	 	paciente_id INT NOT NULL REFERENCES pacientes(id) ON DELETE CASCADE,
		doctor_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		fecha DATE NOT NULL,
		hora VARCHAR(10) NOT NULL,
		motivo TEXT,
		estado VARCHAR(50) NOT NULL DEFAULT 'pendiente' CHECK (estado IN('pendiente', 'confirmada', 'cancelada', 'completada')),
		observaciones TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = DB.Exec(crearTabla)
	if err != nil {
		log.Fatal("❌ Error al crear tabla:", err)
	}

	fmt.Println("✅ Tablas listas")

}

func Close() {
	DB.Close()
	fmt.Println("🔌 Desconectado de PostgreSQL")
}
