package gpt

import (
	"context"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

type OpenAIContextItem struct {
	Message openai.ChatCompletionMessage
	Summary string
	Date    string
}

// NewRequest function creates APIContext instance.
func NewRequest(ctx context.Context, apiKey string, messages []openai.ChatCompletionMessage) (*OpenAIContextItem, error) {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo16K0613,
			MaxTokens:   125,
			Temperature: 0.7,
			Messages:    messages,
		},
	)
	if err != nil {
		return nil, err
	}

	return newOpenAIResponse(resp.Choices[0].Message)
}

func newOpenAIResponse(response openai.ChatCompletionMessage) (*OpenAIContextItem, error) {
	gptContext := &OpenAIContextItem{
		Date:    time.Now().UTC().Format(time.RFC822),
		Message: response,
	}

	var split bool
	gptContext.Message.Content, gptContext.Summary, split = strings.Cut(response.Content, "\n")
	gptContext.Message.Content = strings.TrimSpace(gptContext.Message.Content)
	if !split {
		gptContext.Summary = gptContext.Message.Content
	} else {
		gptContext.Summary = strings.TrimSpace(gptContext.Summary)
	}

	return gptContext, nil
}
