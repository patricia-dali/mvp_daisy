package cadastro

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func ShowCadastroPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/templates/cadastro.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func Cadastro(w http.ResponseWriter, r *http.Request, db *sql.DB) error {
	isAdmin := false
	usernameInput := r.FormValue("username")
	passwordInput := r.FormValue("password")
	emailInput := r.FormValue("email")
	salt := generateSalt()

	passwordWithSalt := []byte(passwordInput + salt)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordWithSalt, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users (username, password, email, salt, admin) VALUES ($1, $2, $3, $4, $5)", usernameInput, hashedPassword, emailInput, salt, isAdmin)

	if err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func generateSalt() string {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", salt)
}
