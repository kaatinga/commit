package gitlet

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaatinga/commit/internal/settings"

	"github.com/go-git/go-git/v5"
)

var Repo *git.Repository

func Init() {
	var err error
	gitPath := getRootRepoFolder(settings.Path)
	if gitPath == "" {
		fmt.Println("Unable to find git repository")
		os.Exit(1)
	}
	Repo, err = git.PlainOpen(gitPath)
	if err != nil {
		fmt.Println("Unable to open git repository")
	}
}

// getRootRepoFolder returns the path to a folder with .git folder inside recursively moving up the folder tree.
func getRootRepoFolder(dir string) string {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println("Unable to get absolute path")
		os.Exit(1)
	}
	pathItems := strings.Split(absDir, string(filepath.Separator))

	fmt.Printf("pathItems: %v\n", pathItems)

	var prefix string
	if filepath.Separator == '/' {
		prefix = "/"
	} else {
		prefix = pathItems[0] + string(filepath.Separator)
	}

	for i := len(pathItems) - 1; i > 1; i-- {
		dirUp := filepath.Join(pathItems[:i+1]...)
		fmt.Printf("dirUp: %s\n", dirUp)
		dirUp = filepath.Join(prefix, dirUp)

		// check if dirUp has .git folder
		if _, err := os.Stat(filepath.Join(dirUp, ".git")); err == nil {
			return dirUp
		}
	}

	return ""
}
