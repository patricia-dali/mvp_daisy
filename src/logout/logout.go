package logout

import (
	"net/http"

	"github.com/gorilla/sessions"
)

func Logout(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["username"] = ""
	session.Values["isAdmin"] = false

	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
