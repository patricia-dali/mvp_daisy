package index

import (
	"net/http"
	"text/template"
	"time"

	openai "exemplo.com/openAI"
	"github.com/gorilla/sessions"
)

func ShowIndexPage(w http.ResponseWriter, r *http.Request) {
	store := sessions.NewCookieStore([]byte("chave-secreta"))
	pergunta := r.FormValue("pergunta")
	parametro := "Vou te passar o nome da minha tabela, as colunas e quero que voce mostre como seria uma consulta sql nela, mas sem quebra de linha somente o texto sem formatação e somente a consulta SQL, sem informações adicionais ou frases que não seja a consulta SQL. nome da tabela: mercado, colunas: nome, preco, quantidade e vendido."
	var respostaAI string
	var tempoDeResposta time.Duration
	var err error

	if pergunta != "" {
		respostaAI, tempoDeResposta, err = openai.OpenAIResponse(pergunta, parametro)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		respostaAI = "Sem pergunta fornecida."
	}

	session, err := store.Get(r, "user-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	isAdmin, isAdminOK := session.Values["isAdmin"].(bool)
	username, usernameOK := session.Values["username"].(string)
	email, emailOK := session.Values["email"].(string)

	if !isAdminOK || !usernameOK || !emailOK {
		http.Error(w, "Erro ao obter dados da sessão", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("./assets/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Username        string
		Email           string
		IsAdmin         bool
		RespostaAI      string
		TempoDeResposta time.Duration
	}{
		Username:        username,
		Email:           email,
		IsAdmin:         isAdmin,
		RespostaAI:      respostaAI,
		TempoDeResposta: tempoDeResposta,
	}

	tmpl.Execute(w, data)
}
