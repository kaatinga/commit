package gpt

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

// NewRequest function creates APIContext instance.
func NewRequest(ctx context.Context, apiKey, content string) (string, error) {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   25,
			Temperature: 0.7,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
