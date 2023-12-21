package settings

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ContextFolder  = "commit"
	KaatingaFolder = ".kaatinga"
	contextFile    = "context.csv"
)

var (
	ContextAbsolutePath string
	RepositoryPath      string
)

func Init() {
	if err := DefinePaths(); err != nil {
		log.Fatal(err)
	}
}

func DefinePaths() error {
	RepositoryPath = getRootRepoFolder(Path)
	if RepositoryPath == "" {
		return fmt.Errorf("unable to find git repository in %s", Path)
	}

	ContextAbsolutePath = filepath.Join(RepositoryPath, KaatingaFolder, ContextFolder, contextFile)
	return nil
}

// getRootRepoFolder returns the path to a folder with .git folder inside recursively moving up the folder tree.
func getRootRepoFolder(dir string) string {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println("Unable to get absolute path to the directory")
	}

	pathItems := strings.Split(absDir, string(filepath.Separator))

	var prefix string
	if filepath.Separator == '/' {
		prefix = "/"
	} else {
		prefix = pathItems[0] + string(filepath.Separator)
	}

	for i := len(pathItems) - 1; i > 1; i-- {
		dirUp := filepath.Join(pathItems[:i+1]...)
		dirUp = filepath.Join(prefix, dirUp)
		if fileInfo, err := os.Stat(filepath.Join(dirUp, ".git")); err == nil {
			if fileInfo.IsDir() {
				return dirUp
			}
		}
	}

	return ""
}
