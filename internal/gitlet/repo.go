package gitlet

import (
	"fmt"

	"github.com/kaatinga/commit/internal/settings"

	"github.com/go-git/go-git/v5"
)

var Repo *git.Repository

func Init() {
	var err error
	Repo, err = git.PlainOpen(settings.RepositoryPath)
	if err != nil {
		fmt.Println("Unable to open git repository")
	}
}
