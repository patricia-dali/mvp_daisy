package financeiro

import (
	"database/sql"
	"net/http"
	"text/template"
)

func ShowFinanceiroPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	tmpl, err := template.ParseFiles("./assets/templates/financeiro.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
