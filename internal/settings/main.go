package settings

import "os"

var (
	RepositoryPath string
	Path           = "."
	APIKey         = os.Getenv("MISTRAL_API_KEY")
	DryRun         = os.Getenv("COMMIT_DRYRUN") == "true"
)
