package gitlet

import (
	"bytes"
	"fmt"
	"github.com/kaatinga/commit/internal/settings"
	"os"
)

const defaultGlobalGitIgnore = "~/.gitignore_global"

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

		fmt.Println("globalGitIgnoreContent", string(globalGitIgnoreContent))

		if !bytes.Contains(globalGitIgnoreContent, []byte(settings.ContextFolder)) {
			globalGitIgnoreMustBeUpdated = true
		}
	}

	// add .commit folder to global .gitignore
	if !globalGitIgnoreMustBeUpdated {
		globalGitIgnoreFile, err := os.OpenFile(defaultGlobalGitIgnore, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
