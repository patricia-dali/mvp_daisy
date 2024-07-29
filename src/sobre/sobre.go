package sobre

import (
	"net/http"
	"text/template"
)

func ShowSobrePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/templates/sobre.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
