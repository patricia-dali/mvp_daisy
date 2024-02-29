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

func envPortOr(port string) string {
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	return ":" + port
}

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

	var port = envPortOr("3000")

	fmt.Printf("Servidor rodando em http://localhost:%s\n", port)
	err = http.ListenAndServe(port, newMux)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

func init() {

	sessionKey = os.Getenv("SESSION_KEY")
}
