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
	parametro := "Eu tenho 4 tabelas no postgresSQL,tabela 1:PRODUTOS com as colunas(ID_PRODUTO, NOME_PRODUTO, CATEGORIA, FABRICANTE_FORNECEDOR, CODIGO_DE_BARRA, STATUS, VALOR_MAIOR, PRECO_MINIMO, ULTIMO_VALOR, MARKUP, QUANTIDADE_VENDIDO, VALOR_VENDIDO), tabela 2:ESTOQUE com as colunas (ID_ESTOQUE, PRODUTO_ID, SALDO_TOTAL, SALDO_RESERVADO, SALDO_DISPONIVEL, PRECO_DE_CUSTO), tabela 3:ESTABELECIMENTOS com as colunas (id_estabelecimento, nome_do_cliente, valor_total, quantidade_total),  tabela 4:PRODUTOS_ESTABELECIMENTOS_RELACIONAMENTO que é uma tabela de relacionamento entre produtos e estabelecimentos com as colunas (ID_RELACIONAMENTO, ID_ESTABELECIMENTO, ID_PRODUTO). Sabendo que estoque é tabela filho de produtos e produtos e estabelecimentos tem uma tabela de relação, me mostre uma consulta sql de acordo com esses dados, sem quebra de linha, somente o comando sql para a pergunta abaixo (coloque todas as letras da consulta em maiusculas, e se não for nome da coluna, sempre faça a busca de algo proximo ao que foi perguntado, não precisa ser exato). NÃO TENHA QUEBRA DE LINHA NO COMANDO!!"
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
		respostaAI, tempoDeResposta, err = openai.OpenAIResponse(pergunta, parametro)
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
		query = "SELECT * FROM produto"
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
