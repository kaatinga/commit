package commit

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/color"
	"github.com/sashabaranov/go-openai"
	"github.com/urfave/cli/v2"

	"github.com/kaatinga/commit/internal/gitlet"
	"github.com/kaatinga/commit/internal/gpt"
	"github.com/kaatinga/commit/internal/settings"
)

const (
	requestTemplate = `
Prepare a commit message for the following changes:

Files:
%s
----------------
Code diff:
%s
`

	apiContext = `System Instruction:

You are the "commit API" with the specific task of generating commit messages for code changes.

For each request, always generate a COMMIT MESSAGE.

COMMIT MESSAGE MUST precisely reflects the provided in the request code diff. Response JSON format and examples:

{
	"message": "Update git pull function to handle the case when a reference is not found"
}

{
	"message": "Add a new function to calculate the sum of two numbers"
}

The response must contain only one field "message" with the commit message.

Required commit message style:
	- commit message must be brief
	- do not prise the changes, we just state what was changed, not why
	- do not put dot in the end of the last sentence

Files are provided by 'git diff --name-only --diff-algorithm=minimal' command.
Code diff is provided by 'git diff --diff-algorithm=minimal' command.
`
)

func Generate(cCtx *cli.Context) error {
	if len(settings.APIKey) != 51 {
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

	commitMessage, err := gpt.NewRequest(ctx, settings.APIKey, gptRequest)
	if err != nil {
		return err
	}

	var gitInfo *gitlet.Message
	gitInfo, err = gitlet.NewGitInfo(commitMessage, config.LocalScope)
	if err != nil {
		return err
	}

	if settings.DryRun {
		fmt.Printf("Dryrun: commit message:\n%s\n", commitMessage)
		return nil
	}

	fmt.Printf("💪 Added commit:\n%s%s%s\n", color.Cyan, commitMessage, color.Reset)
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

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: apiContext,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf(requestTemplate, files, diff),
		},
	}

	return messages, nil
}
