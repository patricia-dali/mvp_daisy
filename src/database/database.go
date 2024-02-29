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

	fmt.Println("Conexão com o banco de dados estabelecida!")
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
