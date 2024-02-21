package index

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
	"time"

	openai "exemplo.com/openAI"
	"github.com/gorilla/sessions"
)

var db *sql.DB

type Produto struct {
	ID         int
	Nome       string
	Preco      float64
	Quantidade int
	Vendido    int
}

func ShowIndexPage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

	if respostaAI != "Sem pergunta fornecida." && respostaAI != "Não tenho resposta para essa pergunta." {
		produtos, err := questionDB(db, respostaAI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, produto := range produtos {
			fmt.Printf("Dados do Produto:\n")
			for col, value := range produto {
				fmt.Printf("%s: %v\n", col, value)
			}
			fmt.Println("----------")
		}
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

func questionDB(db *sql.DB, query string) ([]map[string]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var produtos []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))

		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, err
		}

		produto := make(map[string]interface{})

		for i, col := range columns {
			produto[col] = values[i]
		}

		produtos = append(produtos, produto)
	}

	return produtos, nil
}
