package index

import (
	"database/sql"
	"net/http"
	"strings"
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
	parametro := `O banco de dados possui quatro tabelas:
					1. Tabela: PRODUTOS
					- Colunas:
						- ID_PRODUTO (INTEGER)
						- NOME_PRODUTO (VARCHAR)
						- CATEGORIA (VARCHAR)
						- FABRICANTE_FORNECEDOR (VARCHAR)
						- CODIGO_DE_BARRA (VARCHAR)
						- STATUS (VARCHAR)
						- VALOR_MAIOR (DECIMAL)
						- PRECO_MINIMO (DECIMAL)
						- ULTIMO_VALOR (DECIMAL)
						- MARKUP (DECIMAL)
						- QUANTIDADE_VENDIDO (INTEGER)
						- VALOR_VENDIDO (DECIMAL)

					2. Tabela: ESTOQUE
					- Colunas:
						- ID_ESTOQUE (INTEGER)
						- PRODUTO_ID (INTEGER) - Referência para PRODUTOS(ID_PRODUTO)
						- SALDO_TOTAL (INTEGER)
						- SALDO_RESERVADO (INTEGER)
						- SALDO_DISPONIVEL (INTEGER)
						- PRECO_DE_CUSTO (DECIMAL)

					3. Tabela: ESTABELECIMENTOS
					- Colunas:
						- ID_ESTABELECIMENTO (INTEGER)
						- NOME_DO_CLIENTE (VARCHAR)
						- VALOR_TOTAL (DECIMAL)
						- QUANTIDADE_TOTAL (INTEGER)

					4. Tabela: PRODUTOS_ESTABELECIMENTOS_RELACIONAMENTO
					- Colunas:
						- ID_RELACIONAMENTO (INTEGER)
						- ID_ESTABELECIMENTO (INTEGER) - Referência para ESTABELECIMENTOS(ID_ESTABELECIMENTO)
						- ID_PRODUTO (INTEGER) - Referência para PRODUTOS(ID_PRODUTO)

					Mostre so a consulta SQL, sem quebra de linha e todas as letras maiúsculas.
					Baseado nessa estrutura, aqui está uma pergunta:`

	var respostaAI string
	var tempoDeResposta time.Duration
	var err error

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

	data := struct {
		Username        string
		Email           string
		IsAdmin         bool
		RespostaAI      string
		TempoDeResposta time.Duration
		Produtos        []map[string]interface{}
		Aviso           string
	}{
		Username:        username,
		Email:           email,
		IsAdmin:         isAdmin,
		RespostaAI:      "",
		TempoDeResposta: 0,
		Produtos:        nil,
		Aviso:           "",
	}

	if pergunta != "" {
		respostaAI, tempoDeResposta, err = openai.OpenAIResponse(pergunta, parametro, "", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		forbiddenWords := []string{"DROP", "TRUNCATE", "DELETE"}
		for _, word := range forbiddenWords {
			if strings.Contains(respostaAI, word) {
				data.Aviso = "Operação proibida detectada: " + word
				http.Error(w, "Operação proibida detectada", http.StatusForbidden)
				return
			}
		}

		if respostaAI != "Sem pergunta fornecida." && respostaAI != "Não tenho resposta para essa pergunta." {
			produtos, err := questionDB(db, isAdmin, respostaAI)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data.Produtos = produtos

			respostaAI, tempoDeResposta, err = openai.OpenAIResponse("", "", respostaAI, produtos)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	data.RespostaAI = respostaAI
	data.TempoDeResposta = tempoDeResposta

	tmpl, err := template.ParseFiles("./assets/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func questionDB(db *sql.DB, isAdmin bool, response string) ([]map[string]interface{}, error) {
	var query string

	if isAdmin {
		query = response
	} else {
		query = "SELECT * FROM PRODUTOS"
	}

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
			if dataBytes, ok := values[i].([]uint8); ok {
				produto[col] = string(dataBytes)
			} else {
				produto[col] = values[i]
			}
		}

		produtos = append(produtos, produto)
	}

	return produtos, nil
}
