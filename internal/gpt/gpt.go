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
	return doOpenAIRequest(ctx, openai.NewClient(apiKey), messages, 5, openai.GPT4oMini)
}

func doOpenAIRequest(ctx context.Context, client *openai.Client, messages []openai.ChatCompletionMessage, attempts byte, model string) (*OpenAIContextItem, error) {
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
		if attempts == 0 {
			return nil, err
		}
		var openAIError = new(openai.RequestError)
		if errors.As(err, &openAIError) {
			// if attempts < 3 {
			// 	model = openai.GPT3Dot5Turbo16K0613
			// }
			if openAIError.HTTPStatusCode == 429 || openAIError.HTTPStatusCode >= 500 {
				return doOpenAIRequest(ctx, client, messages, attempts-1, model)
			}
		}
		return nil, err
	}

	output := newOpenAIResponse(resp.Choices[0].Message)
	if output.Message.Content == "" && attempts > 0 {
		fmt.Println("🐌 It will take more time to generate a commit message, please wait...")
		return doOpenAIRequest(ctx, client, messages, attempts-1, model)
	}

	return output, nil
}

func newOpenAIResponse(response openai.ChatCompletionMessage) *OpenAIContextItem {
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

	return gptContext
}
