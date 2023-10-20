package gitlet

import (
	"bytes"
	"errors"
	"github.com/kaatinga/commit/internal/settings"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetDiff() (string, error) {
	return RunCommand(`git diff --diff-algorithm=minimal -- *.go`, settings.Path)
}

func GetFileList() (string, error) {
	return RunCommand(`git diff --name-only --diff-algorithm=minimal`, settings.Path)
}

type GitInfo struct {
	Repo *git.Repository
	Msg  string
	object.Signature
}

func NewGitInfo(msg string) (*GitInfo, error) {
	gitInfo := &GitInfo{Msg: msg}
	var err error
	gitInfo.Repo, err = git.PlainOpen(settings.RepositoryPath)
	if err != nil {
		return nil, err
	}

	// get user info
	var gitConfig *config.Config
	gitConfig, err = gitInfo.Repo.ConfigScoped(config.GlobalScope)
	if err != nil {
		return nil, err
	}
	gitInfo.Name = gitConfig.User.Name
	gitInfo.Email = gitConfig.User.Email
	gitInfo.When = time.Now()

	if gitInfo.Name == "" {
		return nil, errors.New("user name is not set")
	}

	if gitInfo.Email == "" {
		return nil, errors.New("user email is not set")
	}

	return gitInfo, nil
}

// RunCommand executes a command and waits for its output
// specially done because git is messing up STDOUT and STDERR, see this: https://github.com/cli/cli/issues/2984
func RunCommand(cmd string, dir string) (string, error) {
	args := strings.Fields(cmd)
	c := exec.Command(args[0], args[1:]...)
	if dir != "" {
		c.Dir = dir
	}

	stdout, err := c.StdoutPipe()
	if err != nil {
		return "", err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := c.Start(); err != nil {
		return "", err
	}

	output, errStdOut := io.ReadAll(stdout)
	errMsg, errStdErr := io.ReadAll(stderr)

	if err := c.Wait(); err != nil {
		return "", errors.New(string(errMsg))
	}

	if errStdOut != nil {
		return "", errStdOut
	}
	if errStdErr != nil {
		return "", errStdErr
	}

	output = bytes.TrimSuffix(output, []byte{10})

	if len(errMsg) == 0 {
		return string(output), nil
	}

	return string(output), errors.New(string(errMsg))
}
