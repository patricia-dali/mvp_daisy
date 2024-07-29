package login

import (
	"database/sql"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func ShowLoginPage(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if auth, ok := session.Values["username"].(string); ok && auth != "" {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./assets/templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB, store *sessions.CookieStore) {
	usernameInput := r.FormValue("username")
	passwordInput := r.FormValue("password")

	var userID int
	var hashedPassword string
	var email string
	var salt string
	var isAdmin bool

	err := db.QueryRow("SELECT id, password, salt, admin, email FROM users WHERE username = $1", usernameInput).Scan(&userID, &hashedPassword, &salt, &isAdmin, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(passwordInput+salt))
	if err != nil {
		http.Error(w, "Senha inválida", http.StatusUnauthorized)
		return
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["id"] = userID
	session.Values["username"] = usernameInput
	session.Values["isAdmin"] = isAdmin
	session.Values["email"] = email
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/index", http.StatusSeeOther)
}
