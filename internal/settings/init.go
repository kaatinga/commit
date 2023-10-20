package settings

import (
	"log"
	"path/filepath"
)

const (
	ContextFolder = ".commit"
	contextFile   = "context.csv"
)

var (
	RepositoryPath, ContextAbsolutePath string
)

func init() {
	var err error
	RepositoryPath, err = filepath.Abs(Path)
	if err != nil {
		log.Fatal(err)
	}

	ContextAbsolutePath = filepath.Join(RepositoryPath, ContextFolder, contextFile)
}
