package gpt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

// NewRequest function returns commit message.
func NewRequest(ctx context.Context, apiKey string, messages []openai.ChatCompletionMessage) (string, error) {
	return doOpenAIRequest(ctx, openai.NewClient(apiKey), messages, 5, openai.GPT4oMini)
}

type CommitOutput struct {
	Message string `json:"message"`
}

func doOpenAIRequest(ctx context.Context, client *openai.Client, messages []openai.ChatCompletionMessage, attempts byte, model string) (string, error) {
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:               model,
		Messages:            messages,
		MaxCompletionTokens: 135,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})
	if err != nil {
		if attempts == 0 {
			return "", err
		}
		var openAIError = new(openai.RequestError)
		if errors.As(err, &openAIError) {
			if openAIError.HTTPStatusCode == http.StatusTooManyRequests || openAIError.HTTPStatusCode >= 500 {
				return doOpenAIRequest(ctx, client, messages, attempts-1, model)
			}
		}
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("empty response from AI 😢")
	}

	var result CommitOutput
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return "", err
	}

	output := result.Message
	if output == "" {
		if attempts > 0 {
			fmt.Println("🐌 It will take more time to generate a commit message, please wait...")
			return doOpenAIRequest(ctx, client, messages, attempts-1, model)
		}

		return "", errors.New("no response from AI 😢")
	}

	return output, nil
}
