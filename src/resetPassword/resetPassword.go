package resetPassword

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

const tokenSecret = "sua_chave_secreta_para_token"
const resetPasswordURL = "http://localhost:3000/reset-password/token/"

func ShowResetPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./assets/templates/resetPassword.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func SendResetEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email not provided", http.StatusBadRequest)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour).Unix(),
		"email": email,
	})

	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resetLink := resetPasswordURL + tokenString

	err = sendEmail(email, resetLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "flash_message",
		Value: "Um e-mail para redefinir sua senha foi enviado. Verifique sua caixa de entrada.",
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func sendEmail(to, resetLink string) error {
	smtpServer := "smtp.gmail.com"
	smtpPort := 587
	smtpUsername := "redefinirsenha1983@gmail.com"
	smtpPassword := "vtaysanamhtagphi"

	subject := "Redefinição de Senha\n\n"
	body := "Você solicitou a redefinição de senha. Clique no link abaixo para redefinir:\n" + resetLink

	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpServer, smtpPort), auth, smtpUsername, []string{to}, []byte(message))
	if err != nil {
		fmt.Println("Erro ao enviar o e-mail:", err)
		return err
	}

	return nil
}

func ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenString := vars["token"]
	if tokenString == "" {
		http.Error(w, "Token não fornecido", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles("./assets/templates/resetPasswordPage.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {
		newPassword := r.FormValue("newPassword")

		err := updatePassword(token.Claims.(jwt.MapClaims)["email"].(string), newPassword, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, nil)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	vars := mux.Vars(r)
	tokenString := vars["token"]
	if tokenString == "" {
		http.Error(w, "Token não fornecido", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles("./assets/templates/resetPasswordPage.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodPost {

		newPassword := r.FormValue("newPassword")

		err := updatePassword(token.Claims.(jwt.MapClaims)["email"].(string), newPassword, db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, nil)
}

func updatePassword(email, newPassword string, db *sql.DB) error {
	var salt string
	err := db.QueryRow("SELECT salt FROM users WHERE email = $1", email).Scan(&salt)
	if err != nil {
		log.Println("Erro ao recuperar o salt do banco de dados:", err)
		return err
	}

	newPasswordWithSalt := []byte(newPassword + salt)

	hashedNewPassword, err := bcrypt.GenerateFromPassword(newPasswordWithSalt, bcrypt.DefaultCost)
	if err != nil {
		log.Println("Erro ao gerar hash para a nova senha:", err)
		return err
	}

	_, err = db.Exec("UPDATE users SET password = $1 WHERE email = $2", hashedNewPassword, email)
	if err != nil {
		log.Println("Erro ao atualizar a nova senha no banco de dados:", err)
		return err
	}

	return nil
}
