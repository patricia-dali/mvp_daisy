package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	dbHost     = ""
	dbPort     = 0
	dbUser     = ""
	dbPassword = ""
	dbName     = ""
)

const createTableQuery = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(250) NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    salt VARCHAR(250) NOT NULL,
    admin BOOLEAN NOT NULL
);
`

func SetupDatabase() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Erro ao abrir conexão com o banco de dados: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Printf("Erro ao pingar o banco de dados: %v", err)
		return nil, err
	}

	_, err = db.Exec(createTableQuery)
	if err != nil {
		db.Close()
		log.Printf("Erro ao criar a tabela: %v", err)
		return nil, err
	}

	fmt.Println("Conexão com o banco de dados estabelecida e tabela criada!")
	return db, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")

	portStr := os.Getenv("DB_PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			dbPort = port
		} else {
			log.Printf("Erro ao converter a porta do banco de dados para int: %v", err)
		}
	}
}
