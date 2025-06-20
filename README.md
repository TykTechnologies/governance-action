# Governance Action

A reusable Continuous Integration (CI) action written in Go that validates OpenAPI Specifications (OAS) against governance rules. This action is packaged as a Docker container and is compatible with GitHub Actions and GitLab CI.

## Overview

The Governance Action reads an OpenAPI Specification file, sends it to a governance service API along with a rule ID, and reports any governance issues found. It's designed to be used as a quality gate in CI/CD pipelines to ensure APIs meet organizational governance standards.

## Features

- **Multi-platform CI Support**: Works with GitHub Actions and GitLab CI
- **Environment Detection**: Automatically detects the CI environment
- **OpenAPI Specification Validation**: Reads and validates OAS files (JSON/YAML)
- **Governance Rule Evaluation**: Integrates with governance service APIs
- **Detailed Reporting**: Provides clear, actionable feedback on governance issues
- **Docker-based**: Easy deployment and consistent execution environment

## Quick Start

### Prerequisites

- Docker (for running the action)
- Access to a governance service API
- OpenAPI Specification file to validate

### Configuration

The action uses the following input parameters:

| Parameter | Description | Required | Default |
|-----------|-------------|----------|---------|
| `governance_service` | Base URL of the governance service API | Yes* | - |
| `governance_auth` | Authentication token for the governance API | Yes* | - |
| `rule_id` | ID of the governance rule to evaluate against | Yes | - |
| `api_path` | Path to the OpenAPI Specification file | Yes | - |
| `mocked` | Mock mode for testing ("success", "fail", "warning") | No | - |

*Not required when using `mocked` mode for testing.

**Environment Variable Fallbacks:**
The action also supports environment variables:
- `GOVERNANCE_SERVICE` → `governance_service`
- `GOVERNANCE_AUTH` → `governance_auth`
- `RULE_ID` → `rule_id`
- `API_PATH` → `api_path`
- `MOCKED` → `mocked`

## Setup Guides

### GitHub Actions

For detailed GitHub Actions integration instructions, see [GitHub Actions Integration Guide](docs/github-actions-integration.md).

**Quick Example:**
```yaml
- name: Run Governance Check
  uses: tyktechnologies/governance-action@latest
  with:
    governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
    governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
    rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
    api_path: ./api/openapi.yaml
```

**Testing with Mock Mode:**
```yaml
- name: Test Governance Check (Mock Mode)
  uses: tyktechnologies/governance-action@latest
  with:
    rule_id: test-rule-id
    api_path: ./api/openapi.yaml
    mocked: success  # Options: success, fail, warning
```

### GitLab CI

For detailed GitLab CI integration instructions, see [GitLab CI Integration Guide](docs/gitlab-integration.md).

**Quick Example:**
```yaml
governance-check:
  image: ghcr.io/tyktechnologies/governance-action:latest
  variables:
    GOVERNANCE_API_URL: $GOVERNANCE_API_URL
    GOVERNANCE_API_TOKEN: $GOVERNANCE_API_TOKEN
    OAS_FILE_PATH: ./api/openapi.yaml
    RULE_ID: $GOVERNANCE_RULE_ID
  script:
    - /app/governance-action
```

**Testing with Mock Mode:**
```yaml
governance-test:
  image: ghcr.io/tyktechnologies/governance-action:latest
  variables:
    OAS_FILE_PATH: ./api/openapi.yaml
    RULE_ID: test-rule-id
    MOCKED: fail  # Options: success, fail, warning
  script:
    - /app/governance-action
```

### Local Testing

For comprehensive local testing instructions, see [Local Testing Guide](docs/local-testing.md).

**Quick Example:**
```bash
docker run --rm \
  -e GOVERNANCE_API_URL=http://localhost:8080 \
  -e GOVERNANCE_API_TOKEN=your-token \
  -e OAS_FILE_PATH=/workspace/openapi.yaml \
  -e RULE_ID=6853d42c7493327ea805be8a \
  -v $(pwd):/workspace \
  ghcr.io/tyktechnologies/governance-action:latest
```

**Testing with Mock Mode:**
```bash
# Test success scenario
docker run --rm \
  -e OAS_FILE_PATH=/workspace/openapi.yaml \
  -e RULE_ID=test-rule-id \
  -e MOCKED=success \
  -v $(pwd):/workspace \
  ghcr.io/tyktechnologies/governance-action:latest

# Test failure scenario
docker run --rm \
  -e OAS_FILE_PATH=/workspace/openapi.yaml \
  -e RULE_ID=test-rule-id \
  -e MOCKED=fail \
  -v $(pwd):/workspace \
  ghcr.io/tyktechnologies/governance-action:latest

# Test warning scenario
docker run --rm \
  -e OAS_FILE_PATH=/workspace/openapi.yaml \
  -e RULE_ID=test-rule-id \
  -e MOCKED=warning \
  -v $(pwd):/workspace \
  ghcr.io/tyktechnologies/governance-action:latest
```

## Output

The action provides detailed governance analysis reports and sets output variables for use in subsequent CI/CD steps.

### Sample Output

```
================ Governance Analysis Report ================
❌ [ERROR] [paths./test.get.responses.200.content.application/json.schema.properties.message] owasp-string-restricted
    schema of type `string` must specify `format`, `const`, `enum` or `pattern`
    Location: line 1, char 183 - line 1, char 189
    --- OAS snippet ---
       1 | openapi: 3.1.0
    -------------------
❌ [ERROR] [paths./test.get.responses.200.content.application/json.schema.properties.message] owasp-string-limit
    schema of type `string` must specify `maxLength`, `const` or `enum`
    Location: line 1, char 183 - line 1, char 189
    --- OAS snippet ---
       1 | openapi: 3.1.0
    -------------------
⚠️ [WARNING] [paths./test.get.responses] owasp-define-error-responses-500
    missing response code `500` for `GET`
    Location: line 1, char 73 - line 1, char 84
    --- OAS snippet ---
       1 | openapi: 3.1.0
    -------------------
⚠️ [WARNING] [paths./test.get.responses] owasp-define-error-responses-401
    missing response code `401` for `GET`
    Location: line 1, char 73 - line 1, char 84
    --- OAS snippet ---
       1 | openapi: 3.1.0
    -------------------
❌ [ERROR] [paths./test.get.responses.200] owasp-rate-limit
    response with code `200`, must contain one of the defined headers: `{X-RateLimit-Limit} {X-Rate-Limit-Limit} {RateLimit-Limit, RateLimit-Reset} {RateLimit} `
    Location: line 1, char 86 - line 1, char 91
    --- OAS snippet ---
       1 | openapi: 3.1.0
    -------------------
⚠️ [WARNING] [paths./test.get.responses] owasp-define-error-validation
    missing one of `400`, `422`, `4XX` response codes
    Location: line 1, char 73 - line 1, char 84
    --- OAS snippet ---
       1 | openapi: 3.1.0
    -------------------
⚠️ [WARNING] [paths./test.get.responses] owasp-define-error-responses-429
    missing response code `429` for `GET`
    Location: line 1, char 73 - line 1, char 84
    --- OAS snippet ---
       1 | openapi: 3.1.0
    -------------------
===========================================================

Action failed: governance analysis failed with 3 errors and 4 warnings
```

### Output Variables

| Variable | Description |
|----------|-------------|
| `error_count` | Number of governance errors found |
| `warning_count` | Number of governance warnings found |
| `total_issues` | Total number of governance issues found |

## Project Structure

```
governance-action/
├── cmd/
│   └── main.go              # Main application entry point
├── pkg/
│   ├── core/
│   │   └── action.go        # Core action logic
│   └── integrations/
│       ├── governance.go    # Governance API client
│       └── platform.go      # CI platform detection
├── test-data/
│   ├── mock-server.go       # Mock governance service
│   └── openapi.yaml         # Sample OpenAPI spec
├── docs/
│   ├── github-actions-integration.md  # GitHub Actions setup
│   ├── gitlab-integration.md          # GitLab CI setup
│   └── local-testing.md               # Local testing guide
├── Dockerfile               # Multi-stage Docker build
├── action.yml               # GitHub Actions metadata
└── README.md               # This file
```

## Development

### Building

```bash
# Build the Go binary
go build -o main cmd/main.go

# Build the Docker image
docker build -t governance-action .

# Build for multiple platforms
docker buildx build --platform linux/amd64,linux/arm64 -t governance-action .
```

### Testing

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Test locally with mock server
cd test-data && go run mock-server.go &
docker run --rm \
  -e GOVERNANCE_API_URL=http://host.docker.internal:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/test-data/openapi.yaml \
  -e RULE_ID=test-rule-id \
  -v $(pwd):/workspace \
  governance-action
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.