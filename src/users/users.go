package users

import (
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

func ShowUsersPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/templates/users.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
