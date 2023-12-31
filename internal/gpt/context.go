package gpt

import (
	"encoding/base64"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kaatinga/commit/internal/settings"

	"github.com/sashabaranov/go-openai"
)

const csvSeparator = '🧮'

// Persist method saves openai message history to recover context.
func (gptContext *OpenAIContextItem) Persist() error {
	if gptContext == nil {
		return errors.New("empty context cannot be saved")
	}

	contextFile, err := os.OpenFile(settings.ContextAbsolutePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open context file so as to save context: %w", err)
	}
	defer contextFile.Close()

	writer := csv.NewWriter(contextFile)
	writer.Comma = csvSeparator
	defer writer.Flush()

	content := base64.StdEncoding.EncodeToString([]byte(gptContext.Message.Content))
	return writer.Write([]string{gptContext.Date.Format(time.RFC3339), content, gptContext.Message.Role})
}

// OpenContext method reads openai message history to recover context from the file in .commit/context.csv.
func OpenContext() ([]OpenAIContextItem, error) {
	file, err := os.Open(settings.ContextAbsolutePath)
	if err != nil {
		var pathError = new(os.PathError)
		if errors.As(err, &pathError) {
			// create .context folder
			err := os.MkdirAll(filepath.Join(settings.RepositoryPath, settings.KaatingaFolder, settings.ContextFolder), 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create .kaatinga/commit folder: %w", err)
			}

			// create empty context.csv file
			_, err = os.Create(settings.ContextAbsolutePath)
			if err != nil {
				return nil, fmt.Errorf("failed to create context file: %w", err)
			}

			return []OpenAIContextItem{
				{
					Date: time.Now().UTC(),
					Message: openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleUser,
						Content: "[comment] the context history is empty. This is the first request to commit API for this repository",
					},
				},
			}, nil
		}

		return nil, fmt.Errorf("failed to open context file: %w", err)
	}

	var records [][]string
	reader := csv.NewReader(file)
	reader.Comma = csvSeparator
	reader.LazyQuotes = true
	records, err = reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read context file: %w", err)
	}

	var context = make([]OpenAIContextItem, 0, len(records))
	for _, record := range records {
		decodedContent, err := base64.StdEncoding.DecodeString(record[1])
		if err != nil {
			return nil, fmt.Errorf("failed to decode content from context file: %w", err)
		}

		var recordTime time.Time
		recordTime, err = time.Parse(time.RFC3339, record[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse time from context file: %w", err)
		}

		chatCompletion := OpenAIContextItem{
			Date: recordTime,
			Message: openai.ChatCompletionMessage{
				Role: record[2],
			},
		}

		if chatCompletion.Message.Role == openai.ChatMessageRoleUser {
			chatCompletion.Message.Content = "[user]: " + string(decodedContent)
		} else {
			chatCompletion.Message.Content = "[openAI]: " + string(decodedContent)
		}

		context = append(context, chatCompletion)
	}

	return context, nil
}

// func getLatestCommits(path string, i int) ([]OpenAIContextItem, error) {
// get the repository url
// var repoURL string
// repoURL, err = gitlet.RunCommand("git config --get remote.origin.url", "")
// if err != nil {
// 	return nil, err
// }

// get 10 latest commits
// var fromCommits []OpenAIContextItem
// fromCommits, err = getLatestCommits(repoURL, 10)
// if err != nil {
// 	return nil, err
// }
//
// output := append(fromCommits, context...)
//
// // sort by date
// sort.Slice(output, func(i, j int) bool {
// 	return output[i].Date < output[j].Date
// })

// 	gitlet.RunCommand("git clone "+path+" .commit", "")
//
// 	var commits []OpenAIContextItem
// 	for i > 0 {
// 		commits = append(commits, OpenAIContextItem{
// 			Date: commit.Author.When.Format(time.RFC822),
// 			Message: openai.ChatCompletionMessage{
// 				Content: "[user]" + commit.Message,
// 				Role:    openai.ChatMessageRoleUser,
// 			},
// 		})
// 		commit, err = commit.Parents().Next()
// 		if err != nil {
// 			break
// 		}
// 		i--
// 	}
//
// 	return commits, nil
// }
