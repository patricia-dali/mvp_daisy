package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"exemplo.com/database"
	"exemplo.com/router"
	"github.com/gorilla/sessions"
)

var sessionKey = ""

func main() {
	db, err := database.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	newMux := http.NewServeMux()

	store := sessions.NewCookieStore([]byte(sessionKey))

	router.HandleRoutes(newMux, db, store)

	newMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	SERVER_PORT := ":3000"

	if os.Getenv("PORT") != "" {
		SERVER_PORT = "0.0.0.0:" + os.Getenv("PORT")
	}

	fmt.Printf("Servidor rodando em http://localhost:%s\n", SERVER_PORT)
	err = http.ListenAndServe((SERVER_PORT), newMux)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

func init() {

	sessionKey = os.Getenv("SESSION_KEY")
}
