package como

import (
	"net/http"
	"text/template"
)

func ShowComoPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/templates/como.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
