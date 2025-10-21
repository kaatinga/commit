package gitlet

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/kaatinga/commit/internal/settings"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetDiff() (string, error) {
	return RunCommand(`git diff --diff-algorithm=minimal`, settings.RepositoryPath)
}

func GetFileList() (string, error) {
	return RunCommand(`git diff --name-only --diff-algorithm=minimal`, settings.RepositoryPath)
}

type Message struct {
	object.Signature
	Msg string
}

func NewGitInfo(msg string, scope config.Scope) (*Message, error) {
	gitInfo := &Message{Msg: msg}

	// get user info
	gitConfig, err := Repo.ConfigScoped(scope)
	if err != nil {
		return nil, fmt.Errorf("unable to get git config: %w", err)
	}

	gitInfo.Name = gitConfig.User.Name
	gitInfo.Email = gitConfig.User.Email
	gitInfo.When = time.Now()

	if gitInfo.Name == "" {
		if scope == config.GlobalScope {
			return nil, errors.New("user name is not set")
		}

		return NewGitInfo(msg, config.GlobalScope)
	}

	if gitInfo.Email == "" {
		if scope == config.GlobalScope {
			return nil, errors.New("user email is not set")
		}

		return NewGitInfo(msg, config.GlobalScope)
	}

	return gitInfo, nil
}

func RunCommand(cmd string, dir string) (string, error) {
	args := strings.Fields(cmd)
	c := exec.Command(args[0], args[1:]...)
	if dir != "" {
		c.Dir = dir
	}

	output, err := c.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to run command: %w", err)
	}

	return string(output), nil
}
