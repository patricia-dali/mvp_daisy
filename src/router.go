package router

import (
	"database/sql"
	"net/http"

	"exemplo.com/cadastro"
	"exemplo.com/como"
	"exemplo.com/home"
	"exemplo.com/index"
	"exemplo.com/login"
	"exemplo.com/logout"
	"exemplo.com/resetPassword"
	"exemplo.com/sobre"
	"exemplo.com/users"
	"github.com/gorilla/sessions"
)

func AdminAuthMiddleware(next http.HandlerFunc, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "user-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		isAdmin, ok := session.Values["isAdmin"].(bool)
		if !ok || !isAdmin {
			http.Error(w, "Acesso não autorizado", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

func AuthMiddleware(next http.HandlerFunc, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "user-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
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

	indexHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index.ShowIndexPage(w, r, db)
	})
	mux.Handle("/index", AuthMiddleware(indexHandlerFunc, store))

	mux.HandleFunc("/cadastro", cadastro.ShowCadastroPage)

	mux.HandleFunc("/home", home.ShowHomePage)

	mux.HandleFunc("/sobre", sobre.ShowSobrePage)

	mux.HandleFunc("/como", como.ShowComoPage)

	usersHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users.ShowUsersPage(w, r, db)
	})
	mux.Handle("/users", AdminAuthMiddleware(usersHandlerFunc, store))

	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		users.DeleteUserHandler(w, r, db)
	})
	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		users.UpdateUserHandler(w, r, db)
	})

	mux.HandleFunc("/reset-password", resetPassword.ShowResetPage)
	mux.HandleFunc("/send-reset-email", resetPassword.SendResetEmailHandler)

	mux.HandleFunc("/reset-password/token", resetPassword.ResetPasswordPage)

	mux.HandleFunc("/reset-password/submit", func(w http.ResponseWriter, r *http.Request) {
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
