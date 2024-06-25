package openai

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

func OpenAIResponse(pergunta, parametro string, respostaAnterior string, produtos []map[string]interface{}) (string, time.Duration, error) {
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
		for _, produto := range produtos {
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
