package settings

import "os"

var (
	Path   = "."
	APIKey = os.Getenv("OPENAI_API_KEY")
	DryRun = os.Getenv("COMMIT_DRYRUN") == "true"
)
