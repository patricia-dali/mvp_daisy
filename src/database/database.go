package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	dbHost     = ""
	dbPort     = 5432
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
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	fmt.Println("Conexão com o banco de dados estabelecida!")
	return db, nil
}

func init() {
	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
}
