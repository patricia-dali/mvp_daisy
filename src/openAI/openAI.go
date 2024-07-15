package openai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

func PrimeiraPergunta(pergunta string) (string, error) {
	client := openai.NewClient(os.Getenv("CHAVE_OPENAI"))

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: ` O banco de dados possui quatro tabelas:
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
					
					De acordo com as informações acima farei uma pergunta e se a pergunta for sobre o banco de dados, sobre produtos, preço de coisas ou qualquer coisa relacionada ao banco a resposta será 1, se não tiver relação com o banco a resposta será 2.
					Responda somente os números 1 ou 2, nada mais.
						`,
		},
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: pergunta,
	})

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	if err != nil {
		return "", fmt.Errorf("erro ao obter resposta do OpenAI: %v", err)
	}
	return resp.Choices[len(resp.Choices)-1].Message.Content, nil
}

func ResponseAleatoria(pergunta string, username string) (string, time.Duration, error) {
	startTime := time.Now()
	client := openai.NewClient(os.Getenv("CHAVE_OPENAI"))

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Você é um assistente virtual chamado Daysi que ajuda o " + username + " em sua projeto.",
		},
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: pergunta,
	})

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	tempoDeResposta := time.Since(startTime)

	if err != nil {
		return "", 0, fmt.Errorf("erro ao obter resposta do OpenAI: %v", err)
	}

	return resp.Choices[len(resp.Choices)-1].Message.Content, tempoDeResposta, nil
}

func ResponseBD(pergunta, parametro string, respostaAnterior string, results [][]interface{}) (string, time.Duration, error) {
	startTime := time.Now()

	client := openai.NewClient(os.Getenv("CHAVE_OPENAI"))

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "Você é um assistente virtual que ajuda a responder perguntas.",
		},
	}

	if respostaAnterior == "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: parametro + pergunta,
		})
	} else {
		var produtosStrBuilder strings.Builder
		for _, produto := range results {
			produtosStrBuilder.WriteString(fmt.Sprintf("Produto: %v\n", produto))
		}
		produtosStr := produtosStrBuilder.String()

		messageContent := fmt.Sprintf("%s\nResposta anterior: %s\nProdutos:\n%s\nReformule a resposta de forma mais humanizada.",
			parametro, respostaAnterior, produtosStr)
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: messageContent,
		})
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)

	tempoDeResposta := time.Since(startTime)

	if err != nil {
		return "", 0, fmt.Errorf("erro ao obter resposta do OpenAI: %v", err)
	}
	return resp.Choices[len(resp.Choices)-1].Message.Content, tempoDeResposta, nil
}
