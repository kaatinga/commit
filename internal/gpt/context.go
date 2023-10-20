package gpt

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/kaatinga/commit/internal/settings"
	"os"
	"path/filepath"
	"time"

	"github.com/sashabaranov/go-openai"
)

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

	// update .gitignore if needed
	var gitIgnoreFile *os.File
	gitIgnoreFile, err = os.OpenFile(filepath.Join(settings.RepositoryPath, ".gitignore"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open/create .gitignore file: %w", err)
	}
	defer gitIgnoreFile.Close()

	// check that .gitoignore contains .commit folder
	var gitIgnoreContent []byte
	gitIgnoreContent, err = os.ReadFile(filepath.Join(settings.RepositoryPath, ".gitignore"))
	if err != nil {
		return fmt.Errorf("failed to read .gitignore file: %w", err)
	}

	if !bytes.Contains(gitIgnoreContent, []byte(settings.ContextFolder)) {
		_, err = gitIgnoreFile.WriteString(settings.ContextFolder + "/\n")
		if err != nil {
			return fmt.Errorf("failed to write .gitignore file: %w", err)
		}
	}

	writer := csv.NewWriter(contextFile)
	defer writer.Flush()

	return writer.Write([]string{
		gptContext.Date,
		gptContext.Summary,
		gptContext.Message.Role,
	})
}

// OpenContext method reads openai message history to recover context from the file in .commit/context.csv.
func OpenContext() ([]OpenAIContextItem, error) {
	absoluteFileName, err := filepath.Abs(settings.ContextAbsolutePath)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(absoluteFileName)
	if err != nil {
		var pathError = new(os.PathError)
		if errors.As(err, &pathError) {
			// create .context folder
			err := os.Mkdir(filepath.Join(settings.RepositoryPath, settings.ContextFolder), 0755)
			if err != nil {
				return nil, fmt.Errorf("failed to create .commit folder: %w", err)
			}

			// create empty context.csv file
			_, err = os.Create(absoluteFileName)
			if err != nil {
				return nil, fmt.Errorf("failed to create context file: %w", err)
			}

			return []OpenAIContextItem{
				{
					Date: time.Now().UTC().Format(time.RFC822),
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
	records, err = reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var context = make([]OpenAIContextItem, 0, len(records))
	for _, record := range records {
		var content string
		if record[2] == openai.ChatMessageRoleUser {
			content = "[summary]" + record[1]
		} else {
			content = "[openAI]" + record[1]
		}
		context = append(context, OpenAIContextItem{
			Date: record[0],
			Message: openai.ChatCompletionMessage{
				Content: content,
				Role:    record[2],
			},
		})
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