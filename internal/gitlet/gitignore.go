package gitlet

import (
	"bytes"
	"fmt"
	"github.com/kaatinga/commit/internal/settings"
	"os"
	"path/filepath"
)

func UpdateGitIgnore() error {
	// update .gitignore if needed
	gitIgnoreFile, err := os.OpenFile(filepath.Join(settings.RepositoryPath, ".gitignore"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open/create .gitignore file: %w", err)
	}
	defer gitIgnoreFile.Close()

	// check that .gitignore contains .commit folder
	var gitIgnoreContent []byte
	gitIgnoreContent, err = os.ReadFile(filepath.Join(settings.RepositoryPath, ".gitignore"))
	if err != nil {
		return fmt.Errorf("failed to read .gitignore file: %w", err)
	}

	if !bytes.Contains(gitIgnoreContent, []byte(settings.ContextFolder)) {
		_, err = gitIgnoreFile.WriteString(settings.ContextFolder + "/\n")
		if err != nil {
			return fmt.Errorf("failed to write .gitignore file: %w", err)
		}
	}
	return nil
}
