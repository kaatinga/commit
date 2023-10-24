package gitlet

import (
	"bytes"
	"fmt"
	"github.com/kaatinga/commit/internal/settings"
	"os"
	"path/filepath"
)

const defaultGlobalGitIgnoreFile = ".gitignore_global"

func UpdateGitIgnore() error {
	// first, check that global .gitignore exists
	// git config --get core.excludesfile
	var globalGitIgnoreMustBeUpdated bool
	globalGitIgnorePath, err := RunCommand("git config --get core.excludesfile", "")
	fmt.Println("globalGitIgnorePath", globalGitIgnorePath)
	if err != nil {
		return fmt.Errorf("failed to get global .gitignore path: %w", err)
	} else {
		// check that global .gitignore contains .commit folder
		var globalGitIgnoreContent []byte
		globalGitIgnoreContent, err = os.ReadFile(globalGitIgnorePath)
		if err != nil {
			return fmt.Errorf("failed to read global .gitignore file: %w", err)
		}

		fmt.Println("globalGitIgnoreContent:\n", string(globalGitIgnoreContent))

		if !bytes.Contains(globalGitIgnoreContent, []byte(settings.ContextFolder)) {
			globalGitIgnoreMustBeUpdated = true
		}
	}

	// add .commit folder to global .gitignore
	if globalGitIgnoreMustBeUpdated {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		globalGitIgnoreFile, err := os.OpenFile(filepath.Join(homeDir, defaultGlobalGitIgnoreFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open/create global .gitignore file: %w", err)
		}
		defer globalGitIgnoreFile.Close()

		_, err = globalGitIgnoreFile.WriteString(settings.ContextFolder + "/\n")
		if err != nil {
			return fmt.Errorf("failed to write global .gitignore file: %w", err)
		}

		fmt.Printf("📝 Added %s to global .gitignore\n", settings.ContextFolder)
	}

	return nil
}
