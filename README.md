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

The action uses the following environment variables:

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `GOVERNANCE_API_URL` | Base URL of the governance service API | Yes | - |
| `GOVERNANCE_API_TOKEN` | Authentication token for the governance API | Yes | - |
| `OAS_FILE_PATH` | Path to the OpenAPI Specification file | Yes | - |
| `RULE_ID` | ID of the governance rule to evaluate against | Yes | - |
| `VERBOSE` | Enable verbose logging (true/false) | No | false |

## Setup Guides

### GitHub Actions

For detailed GitHub Actions integration instructions, see [GitHub Actions Integration Guide](docs/github-actions-integration.md).

**Quick Example:**
```yaml
- name: Run Governance Check
  uses: docker://ghcr.io/tyktechnologies/governance-action:latest
  env:
    GOVERNANCE_API_URL: ${{ secrets.GOVERNANCE_API_URL }}
    GOVERNANCE_API_TOKEN: ${{ secrets.GOVERNANCE_API_TOKEN }}
    OAS_FILE_PATH: ./api/openapi.yaml
    RULE_ID: ${{ secrets.GOVERNANCE_RULE_ID }}
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

## Output

The action provides detailed governance analysis reports and sets output variables for use in subsequent CI/CD steps.

### Sample Output

```
=== GOVERNANCE ISSUES SUMMARY ===
Found 2 governance issues:

1. [ERROR] Missing API Version Header
   Rule: api-version-header
   Path: /users
   Message: API endpoints should include version header
   Location: paths./users.get

2. [WARNING] Missing Rate Limiting
   Rule: rate-limiting
   Path: /users/{id}
   Message: Consider adding rate limiting to this endpoint
   Location: paths./users/{id}.get

=== RECOMMENDATIONS ===
- Add API version header to all endpoints
- Implement rate limiting for better API protection
- Review security headers configuration

Action failed: governance issues detected
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