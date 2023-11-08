package gpt

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type OpenAIContextItem struct {
	Message openai.ChatCompletionMessage
	Summary string
	Date    time.Time
}

// NewRequest function creates APIContext instance.
func NewRequest(ctx context.Context, apiKey string, messages []openai.ChatCompletionMessage) (*OpenAIContextItem, error) {
	client := openai.NewClient(apiKey)
	resp, err := doOpenAIRequest(ctx, client, messages, 5, openai.GPT3Dot5Turbo1106)
	if err != nil {
		return nil, err
	}

	return newOpenAIResponse(resp.Choices[0].Message)
}

func doOpenAIRequest(ctx context.Context, client *openai.Client, messages []openai.ChatCompletionMessage, attempts byte, model string) (openai.ChatCompletionResponse, error) {
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       model,
			MaxTokens:   135,
			Temperature: 0.7,
			Messages:    messages,
		},
	)
	if err != nil {
		fmt.Printf("openAI error: %T\n", err)
		if attempts == 0 {
			return openai.ChatCompletionResponse{}, err
		}
		var openAIError = new(openai.RequestError)
		if errors.As(err, &openAIError) {
			if attempts < 3 {
				model = openai.GPT3Dot5Turbo16K0613
			}
			if openAIError.HTTPStatusCode == 429 || openAIError.HTTPStatusCode >= 500 {
				return doOpenAIRequest(ctx, client, messages, attempts-1, model)
			}
		}
		return openai.ChatCompletionResponse{}, err
	}

	return resp, nil
}

func newOpenAIResponse(response openai.ChatCompletionMessage) (*OpenAIContextItem, error) {
	gptContext := &OpenAIContextItem{
		Date:    time.Now().UTC(),
		Message: response,
	}

	var split bool
	gptContext.Message.Content, gptContext.Summary, split = strings.Cut(response.Content, "\n")
	gptContext.Message.Content, _ = strings.CutPrefix(gptContext.Message.Content, "[openAI]:")
	gptContext.Message.Content, _ = strings.CutPrefix(gptContext.Message.Content, "Commit message:")
	gptContext.Message.Content = strings.TrimSpace(gptContext.Message.Content)
	if !split {
		gptContext.Summary = gptContext.Message.Content
	} else {
		gptContext.Summary = strings.TrimSpace(gptContext.Summary)
	}

	return gptContext, nil
}
