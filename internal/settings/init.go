package settings

import (
	"fmt"
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

func DefinePaths() (err error) {
	RepositoryPath, err = filepath.Abs(Path)
	if err != nil {
		return fmt.Errorf("unable to get absolute path: %w", err)
	}

	ContextAbsolutePath = filepath.Join(RepositoryPath, KaatingaFolder, ContextFolder, contextFile)
	return nil
}
