package router

import (
	"database/sql"
	"net/http"

	"exemplo.com/cadastro"
	"exemplo.com/index"
	"exemplo.com/login"
	"exemplo.com/logout"
	"exemplo.com/resetPassword"
	"github.com/gorilla/sessions"
)

func AuthMiddleware(next http.HandlerFunc, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "user-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if auth, ok := session.Values["username"].(string); !ok || auth == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func HandleRoutes(mux *http.ServeMux, db *sql.DB, store *sessions.CookieStore) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		login.ShowLoginPage(w, r, store)
	})

	indexHandlerFunc := http.HandlerFunc(index.ShowIndexPage)
	mux.Handle("/index", AuthMiddleware(indexHandlerFunc, store))

	mux.HandleFunc("/cadastro", cadastro.ShowCadastroPage)

	mux.HandleFunc("/reset-password", resetPassword.ShowResetPage)
	mux.HandleFunc("/send-reset-email", resetPassword.SendResetEmailHandler)

	mux.HandleFunc("/reset-password/token/{token}", resetPassword.ResetPasswordPage)

	mux.HandleFunc("/reset-password/submit/{token}", func(w http.ResponseWriter, r *http.Request) {
		resetPassword.ResetPasswordHandler(w, r, db)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		login.Login(w, r, db, store)
	})

	mux.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		cadastro.Cadastro(w, r, db)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logout.Logout(w, r, store)
	})
}
