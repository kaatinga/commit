package gitlet

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kaatinga/commit/internal/settings"

	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetDiff() (string, error) {
	var b strings.Builder

	if hasHEAD() {
		tracked, err := RunCommand(`git diff HEAD --diff-algorithm=minimal`, settings.RepositoryPath)
		if err != nil {
			return tracked, err
		}
		b.WriteString(tracked)
	}

	untracked, err := untrackedFiles()
	if err != nil {
		return b.String(), err
	}
	for _, f := range untracked {
		diff, err := diffNoIndex(f)
		if err != nil {
			return b.String(), err
		}
		b.WriteString(diff)
	}

	return b.String(), nil
}

func GetFileList() (string, error) {
	var b strings.Builder

	if hasHEAD() {
		tracked, err := RunCommand(`git diff HEAD --name-only --diff-algorithm=minimal`, settings.RepositoryPath)
		if err != nil {
			return tracked, err
		}
		b.WriteString(tracked)
	}

	untracked, err := untrackedFiles()
	if err != nil {
		return b.String(), err
	}
	for _, f := range untracked {
		b.WriteString(f)
		b.WriteByte('\n')
	}

	return b.String(), nil
}

func hasHEAD() bool {
	_, err := RunCommand(`git rev-parse --verify HEAD`, settings.RepositoryPath)
	return err == nil
}

func untrackedFiles() ([]string, error) {
	out, err := RunCommand(`git ls-files --others --exclude-standard`, settings.RepositoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to list untracked files: %w", err)
	}

	var files []string
	for _, line := range strings.Split(out, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

func diffNoIndex(path string) (string, error) {
	c := exec.Command("git", "diff", "--no-index", "--diff-algorithm=minimal", "--", os.DevNull, path)
	if settings.RepositoryPath != "" {
		c.Dir = settings.RepositoryPath
	}

	output, err := c.CombinedOutput()
	if err != nil {
		if exitErr, ok := errors.AsType[*exec.ExitError](err); ok && exitErr.ExitCode() == 1 {
			return string(output), nil
		}
		return string(output), fmt.Errorf("failed to diff untracked file %q: %w", path, err)
	}

	return string(output), nil
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
