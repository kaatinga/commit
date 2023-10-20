package commit

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5/plumbing/color"
	"github.com/kaatinga/commit/internal/gitlet"
	"github.com/kaatinga/commit/internal/gpt"
	"github.com/kaatinga/commit/internal/settings"
	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"
)

const (
	requestTemplate = `
Prepare a commit message for the following changes. 
Files:
%s
----------------
Code diff:
%s
`

	apiContext = `System Instruction:

You are the "commit API" with the specific task of generating commit messages and summaries for code changes.

For each request, always generate two lines separated by a newline character. Only one \n character is allowed in the response.
First line MUST be a COMMIT MESSAGE. Second line MUST be a SUMMARY of the changes. No need to add any tag or description for the lines.

COMMIT MESSAGE MUST precisely reflects the provided in the request code diff. Examples:
- Update git pull function to handle the case when a reference is not found.
- Add a new function to calculate the sum of two numbers.

SUMMARY will be stored as context for future requests to you for this repository and MUST contain more details about code changes. Examples:
- "Line 22: added a new function to calculate the sum of two numbers."
- "Line 490: use "if err != nil" to check for errors. Line 3: changed the function name from "Sum" to "Add".

If request does not contain code diff, but the file is present in the Files, then COMMIT MESSAGE and SUMMARY MUST contain information about the file. Examples:
- An unknown change in the file ".gitignore".

Expect the request format:
Files:
<list of files>
----------------
Code diff:
<code diff>

Files are provided by 'git diff --name-only --diff-algorithm=minimal' command.
Code diff is provided by 'git diff --diff-algorithm=minimal *.go' command.

You may receive context with these tags:
- [summary]: Summaries of previous changes.
- [openAI]: Commit messages generated by commit API.
- [user]: Commit messages written by the user.
- [comment]: Additional contextual information.

You MUST NOT use these [summary], [openAI] and any other tags in the response.

Make sure to adhere to this two-line format consistently.
`
)

func Generate(cCtx *cli.Context) error {
	if settings.APIKey == "" {
		return cli.Exit("openAI API key is not specified", 1)
	}

	ctx, cancelFunc := context.WithTimeout(cCtx.Context, 25*time.Second)
	defer cancelFunc()

	gptRequest, err := prepareRequest()
	if err != nil {
		if errors.Is(err, errNoError) {
			return nil
		}
		return err
	}

	var response *gpt.OpenAIContextItem
	response, err = gpt.NewRequest(ctx, settings.APIKey, gptRequest)
	if err != nil {
		return err
	}

	var gitInfo *gitlet.GitInfo
	gitInfo, err = gitlet.NewGitInfo(response.Message.Content)
	if err != nil {
		return err
	}

	if settings.DryRun {
		fmt.Printf("Commit message:\n%s\nsummary:\n%s\n", response.Message.Content, response.Summary)
		return nil
	}

	userRequest := gpt.OpenAIContextItem{
		Message: openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: response.Summary,
		},
		Date: response.Date.Add(-2 * time.Second),
	}
	userRequest.Persist()

	err = response.Persist()
	if err != nil {
		return err
	}

	fmt.Printf("💪 Added commit:\n%s%s%s\n", color.Cyan, response.Message.Content, color.Reset)
	return gitInfo.Commit()
}

var errNoError = errors.New("not an error")

func prepareRequest() ([]openai.ChatCompletionMessage, error) {
	files, err := gitlet.GetFileList()
	if err != nil {
		return nil, fmt.Errorf("failed to get file list: %w", err)
	}

	if files == "" {
		fmt.Printf("%sNothing is changed, no commit is needed.%s\n", color.Green, color.Reset)
		return nil, errNoError
	}

	var diff string
	diff, err = gitlet.GetDiff()
	if err != nil {
		return nil, fmt.Errorf("failed to get diff: %w", err)
	}

	var gptContext []gpt.OpenAIContextItem
	gptContext, err = gpt.OpenContext()
	if err != nil {
		return nil, fmt.Errorf("failed to open history for context: %w", err)
	}

	var messages = make([]openai.ChatCompletionMessage, 0, len(gptContext)+2)

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: apiContext,
	})

	for _, item := range gptContext {
		messages = append(messages, item.Message)
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: fmt.Sprintf(requestTemplate, files, diff),
	})

	// for _, message := range messages {
	// 	fmt.Printf("%s%s%s\n", color.Yellow, message.Content, color.Reset)
	// }

	return messages, nil
}
