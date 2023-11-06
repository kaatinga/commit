package gitlet

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/kaatinga/commit/internal/settings"
	"path/filepath"
	"strings"
)

var Repo *git.Repository

func Init() {
	var err error
	Repo, err = git.PlainOpen(settings.RepositoryPath)
	if err != nil {
		fmt.Println("Unable to open git repository")
	}
}

// getGitlabPath returns the path to a folder with .git folder inside recursively moving up the folder tree.
func findRootGitDir(dir string) (repoDir string, notFound bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", true
	}

	// split dir and gitGS into string slices
	pathItems := strings.Split(dir, string(filepath.Separator))
	homeDirItems := strings.Split(homeDir, string(filepath.Separator))

	if !isDirInHomeDir(pathItems, homeDirItems) {
		return "", true
	}

	var prefix string
	if filepath.Separator == '/' {
		prefix = "/"
	} else {
		prefix = pathItems[0] + string(filepath.Separator)
	}

	for i := len(pathItems); i > len(homeDirItems)+1; i-- {
		dirUp := filepath.Join(pathItems[:i]...)
		dirUp = filepath.Join(prefix, dirUp, ".git")
		if _, err := os.Stat(dirUp); err == nil {
			return dirUp, false
		}
	}

	return "", true
}

func isDirInHomeDir(pathItems, gitGSItems []string) bool {
	if len(pathItems) <= len(gitGSItems) {
		return false
	}

	// check that the dir is in the gitgs path
	for index, item := range gitGSItems {
		if item != pathItems[index] {
			return false
		}
	}

	return true
}
