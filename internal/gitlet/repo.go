package gitlet

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/kaatinga/commit/internal/settings"
)

var Repo *git.Repository

func Init() {
	var err error
	Repo, err = git.PlainOpen(settings.RepositoryPath)
	if err != nil {
		fmt.Println("Unable to open git repository")
	}
}
