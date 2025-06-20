package integrations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// GovernanceClient handles communication with the governance service
type GovernanceClient struct {
	baseURL    string
	authToken  string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewGovernanceClient creates a new governance client
func NewGovernanceClient(baseURL, authToken string, logger *zap.Logger) *GovernanceClient {
	return &GovernanceClient{
		baseURL:   baseURL,
		authToken: authToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// LintResult represents a governance analysis result
type LintResult struct {
	Code     string        `json:"code"`
	Path     []string      `json:"path"`
	Message  string        `json:"message"`
	Severity int           `json:"severity"`
	Range    LintRange     `json:"range"`
	Source   string        `json:"source"`
	API      APIReference  `json:"api"`
	Rule     RuleReference `json:"rule"`
}

// LintRange represents the location of an issue in the source file
type LintRange struct {
	Start LintLocation `json:"start"`
	End   LintLocation `json:"end"`
}

// LintLocation represents a position in the source file
type LintLocation struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// APIReference represents a reference to an API
type APIReference struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RuleReference represents a reference to a rule
type RuleReference struct {
	Name string `json:"name"`
}

// AnalyzeOAS analyzes an OpenAPI specification against a specific rule
func (c *GovernanceClient) AnalyzeOAS(ctx context.Context, oasContent, ruleID, filename string) ([]LintResult, error) {
	c.logger.Info("Starting OAS analysis", zap.String("rule_id", ruleID), zap.String("filename", filename))

	// Create the analysis request in the correct format expected by the governance service
	request := map[string]interface{}{
		"ruleSetSelector": map[string]interface{}{
			"id": ruleID,
		},
		"apiContent": map[string]interface{}{
			"name":    filename,
			"content": oasContent,
		},
	}

	// Make the API call
	results, err := c.makeAnalysisRequest(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to make analysis request: %w", err)
	}

	c.logger.Info("Analysis completed", zap.Int("result_count", len(results)))
	return results, nil
}

// makeAnalysisRequest makes the actual HTTP request to the governance service
func (c *GovernanceClient) makeAnalysisRequest(ctx context.Context, request interface{}) ([]LintResult, error) {
	// For now, we'll use the existing /rulesets/evaluate endpoint
	// In a real implementation, you might need a different endpoint for direct file analysis

	// Convert the request to JSON
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/rulesets/evaluate", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", fmt.Sprintf("%s", c.authToken))

	// Make the request
	c.logger.Debug("Making request to governance service", zap.String("url", url))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check response status
	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Governance service returned error",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response_body", string(body)))
		return nil, fmt.Errorf("governance service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var results []LintResult
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return results, nil
}

// Alternative approach: If the governance service doesn't support direct file analysis,
// we might need to implement a different workflow. Here's a placeholder for that:

// AnalyzeOASWithUpload analyzes an OAS file by first uploading it to the governance service
func (c *GovernanceClient) AnalyzeOASWithUpload(ctx context.Context, oasContent, ruleID string) ([]LintResult, error) {
	c.logger.Info("Starting OAS analysis with upload workflow", zap.String("rule_id", ruleID))

	// This would be the workflow if we need to:
	// 1. Upload the OAS file to create a temporary API
	// 2. Run the evaluation against that API
	// 3. Clean up the temporary API

	// For now, this is a placeholder implementation
	// In a real scenario, you would:
	// 1. Call an upload endpoint to create a temporary API
	// 2. Use the existing /rulesets/evaluate endpoint with the temporary API ID
	// 3. Clean up the temporary API after analysis

	return c.AnalyzeOAS(ctx, oasContent, ruleID, "")
}
