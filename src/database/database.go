package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	dbHost     = ""
	dbPort     = 0
	dbUser     = ""
	dbPassword = ""
	dbName     = ""
)

const createTableDadosEstabelecimentos = `
CREATE TABLE IF NOT EXISTS dados_estabelecimento (
    ID_ESTABELECIMENTO SERIAL PRIMARY KEY,
    NOME_DO_CLIENTE VARCHAR(255)
);
`

const createTableUsers = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(250) NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    salt VARCHAR(250) NOT NULL,
    admin BOOLEAN NOT NULL
);
`

const createTableProduto_Estabelecimento_Relacionamentos = `
CREATE TABLE IF NOT EXISTS PRODUTOS_ESTABELECIMENTOS_RELACIONAMENTO (
    ID_RELACIONAMENTO SERIAL PRIMARY KEY,
    ID_ESTABELECIMENTO INTEGER,
    ID_PRODUTO INTEGER,
    FOREIGN KEY (ID_ESTABELECIMENTO) REFERENCES dados_estabelecimento (ID_ESTABELECIMENTO),
    FOREIGN KEY (ID_PRODUTO) REFERENCES PRODUTOS (ID_PRODUTO)
);
`

const createTableProduto = `
CREATE TABLE IF NOT EXISTS PRODUTOS (
    ID_PRODUTO SERIAL PRIMARY KEY,
    NOME_PRODUTO VARCHAR(255),
    CATEGORIA VARCHAR(100),
    FABRICANTE_FORNECEDOR VARCHAR(255),
    CODIGO_DE_BARRA VARCHAR(100) UNIQUE,
    STATUS VARCHAR(100),
    VALOR_MAIOR DECIMAL(10, 2),
    PRECO_MINIMO DECIMAL(10, 2),
    ULTIMO_VALOR DECIMAL(10, 2),
    MARKUP NUMERIC GENERATED ALWAYS AS ((PRECO_MINIMO/VALOR_MAIOR)-1) STORED,
    QUANTIDADE_VENDIDO INTEGER,
    VALOR_VENDIDO NUMERIC GENERATED ALWAYS AS (PRECO_MINIMO*QUANTIDADE_VENDIDO) STORED
);
`

const createTableEstoque = `
CREATE TABLE IF NOT EXISTS ESTOQUE (
    ID_ESTOQUE SERIAL PRIMARY KEY,
    PRODUTO_ID INTEGER,
    SALDO_TOTAL INTEGER,
    SALDO_RESERVADO INTEGER,
    SALDO_DISPONIVEL NUMERIC GENERATED ALWAYS AS (SALDO_TOTAL-SALDO_RESERVADO) STORED,
    PRECO_DE_CUSTO DECIMAL(10, 2),
    FOREIGN KEY (PRODUTO_ID) REFERENCES PRODUTOS (ID_PRODUTO)
);
`

const createTableEstabelecimentos = `
CREATE OR REPLACE VIEW ESTABELECIMENTOS AS
SELECT
    e.ID_ESTABELECIMENTO,
    e.NOME_DO_CLIENTE,
    SUM(COALESCE(p.VALOR_VENDIDO, 0)) AS VALOR_TOTAL,
    SUM(COALESCE(es.SALDO_TOTAL, 0)) AS QUANTIDADE_TOTAL
FROM
    dados_estabelecimento e
LEFT JOIN
    PRODUTOS_ESTABELECIMENTOS_RELACIONAMENTO per ON e.ID_ESTABELECIMENTO = per.ID_ESTABELECIMENTO
LEFT JOIN
    PRODUTOS p ON per.ID_PRODUTO = p.ID_PRODUTO
LEFT JOIN
    ESTOQUE es ON p.ID_PRODUTO = es.PRODUTO_ID
GROUP BY
    e.ID_ESTABELECIMENTO, e.NOME_DO_CLIENTE
;
`

func SetupDatabase() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("Erro ao abrir conexão com o banco de dados: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		log.Printf("Erro ao pingar o banco de dados: %v", err)
		return nil, err
	}

	_, err = db.Exec(createTableDadosEstabelecimentos + createTableUsers + createTableProduto + createTableEstoque + createTableProduto_Estabelecimento_Relacionamentos + createTableEstabelecimentos)
	if err != nil {
		db.Close()
		log.Printf("Erro ao criar a tabela: %v", err)
		return nil, err
	}

	fmt.Println("Conexão com o banco de dados estabelecida e tabela criada!")
	return db, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost = os.Getenv("DB_HOST")
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")

	portStr := os.Getenv("DB_PORT")
	if portStr != "" {
		port, err := strconv.Atoi(portStr)
		if err == nil {
			dbPort = port
		} else {
			log.Printf("Erro ao converter a porta do banco de dados para int: %v", err)
		}
	}
}
