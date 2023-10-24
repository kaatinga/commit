package settings

import (
	"log"
	"path/filepath"
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

func init() {
	err := DefinePaths()
	if err != nil {
		log.Fatal(err)
	}
}

func DefinePaths() error {
	var err error
	RepositoryPath, err = filepath.Abs(Path)
	if err != nil {
		return err
	}

	ContextAbsolutePath = filepath.Join(RepositoryPath, KaatingaFolder, ContextFolder, contextFile)
	return nil
}
