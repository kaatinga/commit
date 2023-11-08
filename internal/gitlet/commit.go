package gitlet

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func (gi *Message) Commit() error {
	wt, err := Repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}
	dir := wt.Filesystem.Root()

	if err = AddAll(dir); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	if _, err = wt.Commit(gi.Msg, &git.CommitOptions{
		Author: &gi.Signature,
	}); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	return nil
}

func AddAll(dir string) error {
	_, err := RunCommand("git add -A", dir)
	// seems to never return output
	return err
}
