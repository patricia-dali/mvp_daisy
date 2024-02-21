package users

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/sessions"
)

type User struct {
	ID       int
	Username string
	Admin    bool
	Email    string
}

var store = sessions.NewCookieStore([]byte("chave-secreta"))

func getLoggedInUserID(r *http.Request) (int, error) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		return 0, err
	}

	userID, ok := session.Values["id"].(int)
	if !ok {
		return 0, fmt.Errorf("ID do usuário não encontrado na sessão.")
	}

	return userID, nil
}

func ShowUsersPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	loggedInUserID, err := getLoggedInUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users, err := getAllUsers(db, loggedInUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./assets/templates/users.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, users)
}

func getAllUsers(db *sql.DB, loggedInUserID int) ([]User, error) {
	rows, err := db.Query("SELECT id, username, admin, email FROM users WHERE id != $1", loggedInUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Admin, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID do usuário não informado.", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "ID de usuário inválido.", http.StatusBadRequest)
		return
	}

	deleteUser := func() error {
		_, err := db.Exec("DELETE FROM users WHERE id = $1 AND admin != true", userID)
		return err
	}

	err = deleteUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodGet {
		// Exibir o formulário de atualização de usuário
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "ID do usuário não informado.", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "ID de usuário inválido.", http.StatusBadRequest)
			return
		}

		user, err := getUserByID(db, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("./assets/templates/update_user.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, user)
	} else if r.Method == http.MethodPost {
		id := r.FormValue("id")
		if id == "" {
			http.Error(w, "ID do usuário não informado.", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "ID de usuário inválido.", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		admin := r.FormValue("admin") == "on"
		email := r.FormValue("email")

		updateUser := func() error {
			_, err := db.Exec("UPDATE users SET username=$1, admin=$2, email=$3 WHERE id=$4", username, admin, email, userID)
			return err
		}

		err = updateUser()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/users", http.StatusSeeOther)
	} else {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func getUserByID(db *sql.DB, userID int) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, admin, email FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Username, &user.Admin, &user.Email)
	return user, err
}
