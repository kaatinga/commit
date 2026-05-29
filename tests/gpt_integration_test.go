package tests

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kaatinga/mistralai-go"

	"github.com/kaatinga/commit/internal/mistralx"
)

// apiKey returns the Mistral API key from the environment, skipping the test when absent.
func apiKey(t *testing.T) string {
	t.Helper()
	key := strings.TrimSpace(os.Getenv("MISTRAL_API_KEY"))
	if key == "" {
		t.Skip("MISTRAL_API_KEY is not set; skipping integration test")
	}
	return key
}

// systemPrompt mirrors the instruction the commit tool sends so the integration
// test exercises a realistic request shape against the real Mistral API.
const systemPrompt = `You are the "commit API". For each request generate a COMMIT MESSAGE that
precisely reflects the provided code diff. Reply with JSON in the form:

{"message": "Update foo to handle bar"}

The response must contain only one field "message". Keep it brief and do not end it with a dot.`

// sampleDiff is a small, self-explanatory change used to prompt the model.
const sampleDiff = `Files:
internal/gpt/gpt.go
----------------
Code diff:
diff --git a/internal/gpt/gpt.go b/internal/gpt/gpt.go
@@
-	"github.com/sashabaranov/go-openai"
+	mistralai "github.com/kaatinga/mistralai-go"
@@
-	return doOpenAIRequest(ctx, openai.NewClient(apiKey), messages, 5, openai.GPT4oMini)
+	client, err := mistralai.NewClient(apiKey)
+	if err != nil {
+		return "", err
+	}
+	defer client.Close()
+	return doRequest(ctx, client, messages, 5)
`

func TestNewRequest_GeneratesCommitMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	key := apiKey(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messages := []mistralai.ChatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: sampleDiff},
	}

	msg, err := mistralx.NewRequest(ctx, key, messages)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	msg = strings.TrimSpace(msg)
	if msg == "" {
		t.Fatal("expected a non-empty commit message")
	}
	if strings.Contains(msg, "{") || strings.Contains(msg, "\"message\"") {
		t.Fatalf("commit message still looks like raw JSON: %q", msg)
	}
	t.Logf("generated commit message: %q", msg)
}

func TestNewRequest_EmptyKey(t *testing.T) {
	// No network involved: an empty key must fail fast via the client constructor.
	_, err := mistralx.NewRequest(context.Background(), "", []mistralai.ChatMessage{
		{Role: "user", Content: "hi"},
	})
	if err == nil {
		t.Fatal("expected an error for an empty API key")
	}
}

func TestChat_JSONFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}
	key := apiKey(t)

	client, err := mistralai.NewClient(key)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.Chat(ctx, mistralai.ChatRequest{
		System: systemPrompt,
		Input:  sampleDiff,
		Format: mistralai.OutputJSON,
	})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}

	var out mistralx.CommitOutput
	if err = resp.JSON(&out); err != nil {
		t.Fatalf("decode JSON response %q: %v", resp.Content, err)
	}
	if strings.TrimSpace(out.Message) == "" {
		t.Fatalf("expected a non-empty message field, got: %q", resp.Content)
	}
	fmt.Printf("Chat JSON commit message: %q\n", out.Message)
}
