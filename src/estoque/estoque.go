package estoque

import (
	"database/sql"
	"net/http"
	"text/template"
)

func ShowEstoquePage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("./assets/templates/estoque.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
