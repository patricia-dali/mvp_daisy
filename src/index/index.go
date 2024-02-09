package index

import (
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

func ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	store := sessions.NewCookieStore([]byte("chave-secreta"))

	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isAdmin, isAdminOK := session.Values["isAdmin"].(bool)
	username, usernameOK := session.Values["username"].(string)
	email, emailOK := session.Values["email"].(string)

	if !isAdminOK || !usernameOK || !emailOK {
		http.Error(w, "Erro ao obter dados da sessão", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./assets/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string
		Email    string
		IsAdmin  bool
	}{
		Username: username,
		Email:    email,
		IsAdmin:  isAdmin,
	}

	tmpl.Execute(w, data)
}
