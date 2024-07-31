package router

import (
	"database/sql"
	"net/http"

	"exemplo.com/cadastro"
	"exemplo.com/como"
	"exemplo.com/config"
	"exemplo.com/estoque"
	"exemplo.com/financeiro"
	"exemplo.com/home"
	"exemplo.com/index"
	"exemplo.com/login"
	"exemplo.com/logout"
	"exemplo.com/perfil"
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
			http.Redirect(w, r, "/login", http.StatusSeeOther)
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

func NoAuthMiddleware(next http.HandlerFunc, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "user-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		username, ok := session.Values["username"].(string)
		if ok && username != "" {
			http.Redirect(w, r, "/index", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func HandleRoutes(mux *http.ServeMux, db *sql.DB, store *sessions.CookieStore) {
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		login.ShowLoginPage(w, r, store)
	})

	indexHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index.ShowIndexPage(w, r, db)
	})
	mux.Handle("/index", AuthMiddleware(indexHandlerFunc, store))

	financeiroHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		financeiro.ShowFinanceiroPage(w, r, db)
	})
	mux.Handle("/financeiro", AuthMiddleware(financeiroHandlerFunc, store))

	perfilHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		perfil.ShowPerfilPage(w, r, db)
	})
	mux.Handle("/perfil", AuthMiddleware(perfilHandlerFunc, store))

	estoqueHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		estoque.ShowEstoquePage(w, r, db)
	})
	mux.Handle("/estoque", AuthMiddleware(estoqueHandlerFunc, store))

	configHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.ShowConfigPage(w, r, db)
	})
	mux.Handle("/config", AuthMiddleware(configHandlerFunc, store))

	mux.HandleFunc("/cadastro", cadastro.ShowCadastroPage)

	mux.HandleFunc("/", NoAuthMiddleware(home.ShowHomePage, store))

	mux.HandleFunc("/sobre", NoAuthMiddleware(sobre.ShowSobrePage, store))

	mux.HandleFunc("/como", NoAuthMiddleware(como.ShowComoPage, store))

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

	mux.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		login.Login(w, r, db, store)
	})

	mux.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		cadastro.Cadastro(w, r, db)
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logout.Logout(w, r, store)
	})
}
