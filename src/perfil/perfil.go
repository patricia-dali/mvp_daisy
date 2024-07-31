package perfil

import (
	"database/sql"
	"net/http"
	"text/template"
)

func ShowPerfilPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("./assets/templates/perfil.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
