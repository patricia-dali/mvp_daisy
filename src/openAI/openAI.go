package openai

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

func OpenAIResponse(pergunta, parametro string) (string, time.Duration, error) {
	startTime := time.Now()

	client := openai.NewClient(os.Getenv("CHAVE_OPENAI"))

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Você é um assistente virtual que ajuda a responder perguntas.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: parametro + pergunta,
				},
			},
		},
	)

	tempoDeResposta := time.Since(startTime)

	if err != nil {
		return "", 0, fmt.Errorf("erro ao obter resposta do OpenAI: %v", err)
	}
	return resp.Choices[len(resp.Choices)-1].Message.Content, tempoDeResposta, nil
}
