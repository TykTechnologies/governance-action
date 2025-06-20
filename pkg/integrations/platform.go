package integrations

import "os"

// DetectCI detects the running CI platform
func DetectCI() string {
	switch {
	case os.Getenv("GITHUB_ACTIONS") == "true":
		return "github"
	case os.Getenv("GITLAB_CI") == "true":
		return "gitlab"
	default:
		return "local"
	}
}

// GetContext extracts context information based on the CI platform
func GetContext(ci string) map[string]string {
	switch ci {
	case "github":
		return map[string]string{
			"repository": os.Getenv("GITHUB_REPOSITORY"),
			"commit":     os.Getenv("GITHUB_SHA"),
			"branch":     os.Getenv("GITHUB_REF_NAME"),
			"actor":      os.Getenv("GITHUB_ACTOR"),
			"workflow":   os.Getenv("GITHUB_WORKFLOW"),
			"run_id":     os.Getenv("GITHUB_RUN_ID"),
		}
	case "gitlab":
		return map[string]string{
			"repository": os.Getenv("CI_PROJECT_PATH"),
			"commit":     os.Getenv("CI_COMMIT_SHA"),
			"branch":     os.Getenv("CI_COMMIT_BRANCH"),
			"actor":      os.Getenv("GITLAB_USER_NAME"),
			"pipeline":   os.Getenv("CI_PIPELINE_ID"),
			"job":        os.Getenv("CI_JOB_ID"),
		}
	default:
		return map[string]string{"env": "local"}
	}
}
