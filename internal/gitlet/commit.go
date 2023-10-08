package gitlet

import "github.com/go-git/go-git/v5"

func (gi *GitInfo) Commit() error {
	repo := gi.Repo

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	dir := wt.Filesystem.Root()

	err = AddAll(dir)
	if err != nil {
		return err
	}

	_, err = wt.Commit(gi.Msg, &git.CommitOptions{
		Author: &gi.Signature,
	})
	return err
}

func AddAll(dir string) error {
	_, err := RunCommand("git add -A", dir)
	// seems to never return output
	return err
}
