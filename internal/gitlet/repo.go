package gitlet

import (
	"fmt"

	"github.com/go-git/go-git/v5"

	"github.com/kaatinga/commit/internal/settings"
)

var Repo *git.Repository

func Open(path string) error {
	var err error
	Repo, err = git.PlainOpenWithOptions(path, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return fmt.Errorf("unable to open git repository in %s: %w", path, err)
	}

	wt, err := Repo.Worktree()
	if err != nil {
		return fmt.Errorf("unable to get worktree: %w", err)
	}
	settings.RepositoryPath = wt.Filesystem.Root()

	return nil
}
