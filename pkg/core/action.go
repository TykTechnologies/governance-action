package core

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TykTechnologies/governance-action/pkg/integrations"
	"go.uber.org/zap"
)

// RunAction is the main entry point for the governance action
func RunAction(logger *zap.Logger) error {
	logger.Info("Starting governance action")

	// Detect CI platform
	ci := integrations.DetectCI()
	logger.Info("Detected CI platform", zap.String("platform", ci))

	// Get context information
	ciContext := integrations.GetContext(ci)
	logger.Info("Retrieved context", zap.Any("context", ciContext))

	// Get configuration from environment
	config, err := getConfiguration()
	if err != nil {
		logger.Error("Failed to get configuration", zap.Error(err))
		return fmt.Errorf("configuration error: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		logger.Error("Invalid configuration", zap.Error(err))
		return fmt.Errorf("invalid configuration: %w", err)
	}

	var results []integrations.LintResult

	// Check if mocked mode is enabled
	if config.Mocked != "" {
		logger.Info("Running in mocked mode", zap.String("mocked_type", config.Mocked))

		// Generate mock results based on the mocked type
		results = generateMockResults(config.Mocked, config.RuleID)
		logger.Info("Generated mock results", zap.Int("result_count", len(results)), zap.String("mocked_type", config.Mocked))
	} else {
		// Normal mode - create governance client and analyze
		client := integrations.NewGovernanceClient(config.GovernanceService, config.GovernanceAuth, logger)

		// Read and validate the OAS file
		oasContent, err := readOASFile(config.APIPath)
		if err != nil {
			logger.Error("Failed to read OAS file", zap.Error(err), zap.String("path", config.APIPath))
			return fmt.Errorf("failed to read OAS file: %w", err)
		}

		// Analyze the OAS file
		results, err = client.AnalyzeOAS(context.Background(), oasContent, config.RuleID)
		if err != nil {
			logger.Error("Failed to analyze OAS", zap.Error(err))
			return fmt.Errorf("failed to analyze OAS: %w", err)
		}
	}

	// Process and report results
	if err := processResults(results, logger); err != nil {
		logger.Error("Failed to process results", zap.Error(err))
		return fmt.Errorf("failed to process results: %w", err)
	}

	logger.Info("Governance action completed successfully")
	return nil
}

// Configuration holds the action configuration
type Configuration struct {
	GovernanceService string
	GovernanceAuth    string
	RuleID            string
	APIPath           string
	Mocked            string
}

// getConfiguration retrieves configuration from environment variables
func getConfiguration() (*Configuration, error) {
	config := &Configuration{
		GovernanceService: os.Getenv("INPUT_GOVERNANCE_SERVICE"),
		GovernanceAuth:    os.Getenv("INPUT_GOVERNANCE_AUTH"),
		RuleID:            os.Getenv("INPUT_RULE_ID"),
		APIPath:           os.Getenv("INPUT_API_PATH"),
		Mocked:            os.Getenv("INPUT_MOCKED"),
	}

	// Fallback to direct environment variables if INPUT_ prefixed ones are not set
	if config.GovernanceService == "" {
		config.GovernanceService = os.Getenv("GOVERNANCE_SERVICE")
	}
	if config.GovernanceAuth == "" {
		config.GovernanceAuth = os.Getenv("GOVERNANCE_AUTH")
	}
	if config.RuleID == "" {
		config.RuleID = os.Getenv("RULE_ID")
	}
	if config.APIPath == "" {
		config.APIPath = os.Getenv("API_PATH")
	}
	if config.Mocked == "" {
		config.Mocked = os.Getenv("MOCKED")
	}

	// GitLab CI specific fallbacks
	if config.GovernanceService == "" {
		config.GovernanceService = os.Getenv("GOVERNANCE_API_URL")
	}
	if config.GovernanceAuth == "" {
		config.GovernanceAuth = os.Getenv("GOVERNANCE_API_TOKEN")
	}
	if config.RuleID == "" {
		config.RuleID = os.Getenv("GOVERNANCE_RULE_ID")
	}
	if config.APIPath == "" {
		config.APIPath = os.Getenv("OAS_FILE_PATH")
	}

	return config, nil
}

// Validate checks if the configuration is valid
func (c *Configuration) Validate() error {
	// If mocked mode is enabled, validate the mocked value
	if c.Mocked != "" {
		if c.Mocked != "success" && c.Mocked != "fail" && c.Mocked != "warning" {
			return fmt.Errorf("mocked must be one of: success, fail, warning")
		}
		// In mocked mode, governance service and auth are not required
		if c.RuleID == "" {
			return fmt.Errorf("rule_id is required")
		}
		if c.APIPath == "" {
			return fmt.Errorf("api_path is required")
		}
		return nil
	}

	// Normal mode validation
	if c.GovernanceService == "" {
		return fmt.Errorf("governance_service is required")
	}
	if c.GovernanceAuth == "" {
		return fmt.Errorf("governance_auth is required")
	}
	if c.RuleID == "" {
		return fmt.Errorf("rule_id is required")
	}
	if c.APIPath == "" {
		return fmt.Errorf("api_path is required")
	}
	return nil
}

// readOASFile reads the OAS file from the specified path
func readOASFile(path string) (string, error) {
	// Resolve relative paths
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path: %w", err)
		}
		path = absPath
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return string(content), nil
}

// generateMockResults creates predefined governance analysis results for testing
func generateMockResults(mockedType string, ruleID string) []integrations.LintResult {
	switch mockedType {
	case "success":
		// Return empty results for success
		return []integrations.LintResult{}

	case "warning":
		// Return warning results
		return []integrations.LintResult{
			{
				Code:     "mock-warning-001",
				Path:     []string{"paths", "/test", "get", "responses"},
				Message:  "Mock warning: Consider adding rate limiting headers",
				Severity: 1, // Warning
				Range: integrations.LintRange{
					Start: integrations.LintLocation{Line: 10, Character: 5},
					End:   integrations.LintLocation{Line: 10, Character: 15},
				},
				Source: "mock-source",
				API: integrations.APIReference{
					ID:   "mock-api-id",
					Name: "Mock API",
				},
				Rule: integrations.RuleReference{
					Name: ruleID,
				},
			},
			{
				Code:     "mock-warning-002",
				Path:     []string{"paths", "/test", "get", "security"},
				Message:  "Mock warning: Consider adding authentication",
				Severity: 1, // Warning
				Range: integrations.LintRange{
					Start: integrations.LintLocation{Line: 8, Character: 3},
					End:   integrations.LintLocation{Line: 8, Character: 12},
				},
				Source: "mock-source",
				API: integrations.APIReference{
					ID:   "mock-api-id",
					Name: "Mock API",
				},
				Rule: integrations.RuleReference{
					Name: ruleID,
				},
			},
		}

	case "fail":
		// Return error results
		return []integrations.LintResult{
			{
				Code:     "mock-error-001",
				Path:     []string{"paths", "/test", "get", "responses"},
				Message:  "Mock error: Missing required 401 response code",
				Severity: 0, // Error
				Range: integrations.LintRange{
					Start: integrations.LintLocation{Line: 10, Character: 5},
					End:   integrations.LintLocation{Line: 10, Character: 15},
				},
				Source: "mock-source",
				API: integrations.APIReference{
					ID:   "mock-api-id",
					Name: "Mock API",
				},
				Rule: integrations.RuleReference{
					Name: ruleID,
				},
			},
			{
				Code:     "mock-error-002",
				Path:     []string{"paths", "/test", "get", "responses", "200"},
				Message:  "Mock error: Missing rate limiting headers in 200 response",
				Severity: 0, // Error
				Range: integrations.LintRange{
					Start: integrations.LintLocation{Line: 12, Character: 7},
					End:   integrations.LintLocation{Line: 12, Character: 10},
				},
				Source: "mock-source",
				API: integrations.APIReference{
					ID:   "mock-api-id",
					Name: "Mock API",
				},
				Rule: integrations.RuleReference{
					Name: ruleID,
				},
			},
			{
				Code:     "mock-warning-003",
				Path:     []string{"paths", "/test", "get", "security"},
				Message:  "Mock warning: Consider adding authentication",
				Severity: 1, // Warning
				Range: integrations.LintRange{
					Start: integrations.LintLocation{Line: 8, Character: 3},
					End:   integrations.LintLocation{Line: 8, Character: 12},
				},
				Source: "mock-source",
				API: integrations.APIReference{
					ID:   "mock-api-id",
					Name: "Mock API",
				},
				Rule: integrations.RuleReference{
					Name: ruleID,
				},
			},
		}

	default:
		return []integrations.LintResult{}
	}
}

// processResults handles the analysis results and determines success/failure
func processResults(results []integrations.LintResult, logger *zap.Logger) error {
	if len(results) == 0 {
		logger.Info("No governance issues found")
		return nil
	}

	// Read OAS file lines for snippet printing
	oasLines := []string{}
	apiPath := os.Getenv("INPUT_API_PATH")
	if apiPath == "" {
		apiPath = os.Getenv("API_PATH")
	}
	if apiPath != "" {
		if file, err := os.Open(apiPath); err == nil {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				oasLines = append(oasLines, scanner.Text())
			}
			file.Close()
		}
	}

	fmt.Println("\n================ Governance Analysis Report ================")
	errorCount := 0
	warningCount := 0
	for _, result := range results {
		sev := "INFO"
		icon := "ℹ️"
		switch result.Severity {
		case 0:
			sev = "ERROR"
			icon = "❌"
			errorCount++
		case 1:
			sev = "WARNING"
			icon = "⚠️"
			warningCount++
		}
		path := strings.Join(result.Path, ".")
		fmt.Printf("%s [%s] [%s] %s\n    %s\n    Location: line %d, char %d - line %d, char %d\n",
			icon, sev, path, result.Rule.Name, result.Message,
			result.Range.Start.Line, result.Range.Start.Character,
			result.Range.End.Line, result.Range.End.Character)

		// Print OAS snippet if available
		if len(oasLines) > 0 && int(result.Range.Start.Line) > 0 && int(result.Range.End.Line) <= len(oasLines) {
			fmt.Println("    --- OAS snippet ---")
			for i := int(result.Range.Start.Line) - 1; i < int(result.Range.End.Line) && i < len(oasLines); i++ {
				fmt.Printf("    %4d | %s\n", i+1, oasLines[i])
			}
			fmt.Println("    -------------------")
		}
	}
	fmt.Println("===========================================================\n")

	// Set output variables for GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		setGitHubOutput("error_count", fmt.Sprintf("%d", errorCount))
		setGitHubOutput("warning_count", fmt.Sprintf("%d", warningCount))
		setGitHubOutput("total_issues", fmt.Sprintf("%d", len(results)))
	}

	// Set output variables for GitLab CI
	if os.Getenv("GITLAB_CI") == "true" {
		setGitLabOutput("error_count", fmt.Sprintf("%d", errorCount))
		setGitLabOutput("warning_count", fmt.Sprintf("%d", warningCount))
		setGitLabOutput("total_issues", fmt.Sprintf("%d", len(results)))
	}

	// Fail if there are errors
	if errorCount > 0 {
		return fmt.Errorf("governance analysis failed with %d errors and %d warnings", errorCount, warningCount)
	}

	return nil
}

// setGitHubOutput sets a GitHub Actions output variable
func setGitHubOutput(name, value string) {
	if outputFile := os.Getenv("GITHUB_OUTPUT"); outputFile != "" {
		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			fmt.Fprintf(f, "%s=%s\n", name, value)
		}
	}
}

// setGitLabOutput sets a GitLab CI output variable
func setGitLabOutput(name, value string) {
	// GitLab CI uses environment variables for outputs
	// We can also write to a file that can be sourced in subsequent jobs
	outputFile := os.Getenv("GITLAB_OUTPUT_FILE")
	if outputFile == "" {
		outputFile = "governance_output.env"
	}

	f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		defer f.Close()
		fmt.Fprintf(f, "export %s=%s\n", name, value)
	}

	// Also set as environment variable for current job
	os.Setenv(name, value)
}
