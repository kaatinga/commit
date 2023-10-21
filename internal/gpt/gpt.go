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
	Date    time.Time
}

// NewRequest function creates APIContext instance.
func NewRequest(ctx context.Context, apiKey string, messages []openai.ChatCompletionMessage) (*OpenAIContextItem, error) {
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo16K0613,
			MaxTokens:   135,
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
		Date:    time.Now().UTC(),
		Message: response,
	}

	var split bool
	gptContext.Message.Content, gptContext.Summary, split = strings.Cut(response.Content, "\n")
	gptContext.Message.Content, _ = strings.CutPrefix(gptContext.Message.Content, "[openAI]:")
	gptContext.Message.Content = strings.TrimSpace(gptContext.Message.Content)
	if !split {
		gptContext.Summary = gptContext.Message.Content
	} else {
		gptContext.Summary = strings.TrimSpace(gptContext.Summary)
	}

	return gptContext, nil
}
