package mistralx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kaatinga/mistralai-go"
)

// model is the Mistral chat model used to generate commit messages.
const model = mistralai.ChatModelMistralSmallLatest

// NewRequest returns a commit message generated from the provided chat messages.
func NewRequest(ctx context.Context, apiKey string, messages []mistralai.ChatMessage) (string, error) {
	client, err := mistralai.NewClient(apiKey)
	if err != nil {
		return "", err
	}
	defer client.Close()

	return doRequest(ctx, client, messages, 5)
}

type CommitOutput struct {
	Message string `json:"message"`
}

// doRequest calls the Mistral chat completions API and parses the commit message
// out of the JSON response. Transient HTTP failures (429/5xx) are retried by the
// client itself; attempts here covers the rare case of an empty model response.
func doRequest(ctx context.Context, client mistralai.Client, messages []mistralai.ChatMessage, attempts byte) (string, error) {
	resp, err := client.ChatCompletion(ctx, mistralai.ChatCompletionRequest{
		Model:     model,
		Messages:  messages,
		MaxTokens: 135,
		ResponseFormat: &mistralai.ResponseFormat{
			Type: "json_object",
		},
	})
	if err != nil {
		return "", err
	}

	content, err := resp.FirstChoiceContent()
	if err != nil {
		return "", errors.New("empty response from AI 😢")
	}

	var result CommitOutput
	if err = json.Unmarshal([]byte(content), &result); err != nil {
		return "", err
	}

	if result.Message == "" {
		if attempts > 0 {
			fmt.Println("🐌 It will take more time to generate a commit message, please wait...")
			return doRequest(ctx, client, messages, attempts-1)
		}

		return "", errors.New("no response from AI 😢")
	}

	return result.Message, nil
}
