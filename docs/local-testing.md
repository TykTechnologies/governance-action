# Local Testing Guide

This document provides comprehensive instructions for testing the Governance Action locally before deploying it to your CI/CD pipelines.

## Overview

Local testing allows you to validate your governance action configuration, test with different OpenAPI specifications, and debug issues before committing changes to your repository.

## Prerequisites

Before testing locally, ensure you have the following installed:

- **Docker**: For running the action in a containerized environment
- **Go 1.21+**: For direct execution and development
- **act** (optional): For testing GitHub Actions workflows locally
- **Sample OpenAPI file**: For testing the action functionality
- **Mock Governance Service**: For simulating governance API responses

## Testing Methods

### Method 1: Using Docker Directly

This is the most common and recommended approach for local testing.

#### Step 1: Build the Docker Image

```bash
# Build the governance action Docker image
docker build -t governance-action .
```

#### Step 2: Create a Test OpenAPI File

If you don't have an OpenAPI file to test with, create a sample one:

```bash
# Create a test OpenAPI file
cat > test-data/sample-openapi.yaml << 'EOF'
openapi: 3.0.0
info:
  title: Sample API
  version: 1.0.0
  description: A sample API for testing governance rules
paths:
  /users:
    get:
      summary: Get users
      responses:
        '200':
          description: List of users
          content:
            application/json:
              schema:
                type: array
                items:
                  type: object
                  properties:
                    id:
                      type: integer
                    name:
                      type: string
  /users/{id}:
    get:
      summary: Get user by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: User details
        '404':
          description: User not found
EOF
```

#### Step 3: Start the Mock Governance Server

The repository includes a mock governance server for testing. Start it in a separate terminal:

```bash
# Navigate to the test-data directory
cd test-data

# Start the mock server
go run mock-server.go
```

#### Step 4: Test the Action

```bash
# Test the action with Docker
docker run --rm \
  -e GOVERNANCE_API_URL=http://host.docker.internal:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/test-data/sample-openapi.yaml \
  -e RULE_ID=test-rule-id \
  -e VERBOSE=true \
  -v $(pwd):/workspace \
  governance-action
```

**Note**: 
- Use `host.docker.internal:8080` on macOS/Windows
- Use `172.17.0.1:8080` on Linux if `host.docker.internal` doesn't work

#### Expected Output

You should see output similar to:

```
2024-01-15T10:30:00.000Z    INFO    Starting governance action
2024-01-15T10:30:00.000Z    INFO    Environment: local
2024-01-15T10:30:00.000Z    INFO    Reading OAS file: /workspace/test-data/sample-openapi.yaml
2024-01-15T10:30:00.000Z    INFO    OAS file read successfully (1024 bytes)
2024-01-15T10:30:00.000Z    INFO    Calling governance API: http://host.docker.internal:8080/api/rulesets/evaluate
2024-01-15T10:30:00.000Z    INFO    Governance API response received
2024-01-15T10:30:00.000Z    WARN    Governance issues found: 2

================ Governance Analysis Report ================
❌ [ERROR] [paths./test.get.responses.200.content.application/json.schema.properties.message] owasp-string-restricted
    schema of type `string` must specify `format`, `const`, `enum` or `pattern`
    Location: line 1, char 183 - line 1, char 189
    --- OAS snippet ---
       1 | openapi: 3.1.0
    -------------------
❌ [ERROR] [paths./test.get.responses.200] owasp-rate-limit
    response with code `200`, must contain one of the defined headers: `{X-RateLimit-Limit} {X-Rate-Limit-Limit} {RateLimit-Limit, RateLimit-Reset} {RateLimit} `
    Location: line 1, char 86 - line 1, char 91
    --- OAS snippet ---
       1 | openapi: 3.1.0
    -------------------
===========================================================

Action failed: governance analysis failed with 2 errors and 0 warnings
```

### Method 2: Using act (GitHub Actions Local Testing)

This method allows you to test your GitHub Actions workflow locally.

#### Step 1: Install act

```bash
# macOS
brew install act

# Linux
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Windows (using Chocolatey)
choco install act-cli
```

#### Step 2: Create a Test Workflow

Create a test workflow file:

```bash
mkdir -p .github/workflows
cat > .github/workflows/test.yml << 'EOF'
name: Test Governance Action
on: [push]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run Governance Check
        uses: ./
        env:
          GOVERNANCE_API_URL: http://host.docker.internal:8080
          GOVERNANCE_API_TOKEN: mock-token
          OAS_FILE_PATH: ./test-data/sample-openapi.yaml
          RULE_ID: test-rule-id
          VERBOSE: true
EOF
```

#### Step 3: Start the Mock Server

In a separate terminal:

```bash
cd test-data
go run mock-server.go
```

#### Step 4: Run the Workflow Locally

```bash
# Run the workflow locally
act push

# Or run a specific job
act push -j test
```

### Method 3: Direct Go Execution (Development)

This method is useful for development and debugging.

#### Step 1: Set Environment Variables

```bash
# Set up environment variables
export GOVERNANCE_API_URL=http://localhost:8080
export GOVERNANCE_API_TOKEN=mock-token
export OAS_FILE_PATH=./test-data/sample-openapi.yaml
export RULE_ID=test-rule-id
export VERBOSE=true
```

#### Step 2: Start the Mock Server

In a separate terminal:

```bash
cd test-data
go run mock-server.go
```

#### Step 3: Run the Action

```bash
# Run the action directly
go run cmd/main.go

# Or build and run
go build -o main cmd/main.go
./main
```

## Testing Different Scenarios

### Testing with Different OAS Files

Test the action with various OpenAPI specifications:

```bash
# Test with JSON format
docker run --rm \
  -e GOVERNANCE_API_URL=http://host.docker.internal:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/test-data/openapi.json \
  -e RULE_ID=test-rule-id \
  -v $(pwd):/workspace \
  governance-action

# Test with different file paths
docker run --rm \
  -e GOVERNANCE_API_URL=http://host.docker.internal:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/api/v1/openapi.yaml \
  -e RULE_ID=test-rule-id \
  -v $(pwd):/workspace \
  governance-action
```

### Testing Error Scenarios

#### Invalid OAS File

```bash
# Create an invalid OAS file
echo "invalid: yaml: content" > test-data/invalid.yaml

# Test with invalid file
docker run --rm \
  -e GOVERNANCE_API_URL=http://host.docker.internal:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/test-data/invalid.yaml \
  -e RULE_ID=test-rule-id \
  -v $(pwd):/workspace \
  governance-action
```

#### Missing File

```bash
# Test with non-existent file
docker run --rm \
  -e GOVERNANCE_API_URL=http://host.docker.internal:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/non-existent.yaml \
  -e RULE_ID=test-rule-id \
  -v $(pwd):/workspace \
  governance-action
```

#### Network Issues

```bash
# Test with unreachable governance service
docker run --rm \
  -e GOVERNANCE_API_URL=http://unreachable:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/test-data/sample-openapi.yaml \
  -e RULE_ID=test-rule-id \
  -v $(pwd):/workspace \
  governance-action
```

### Testing CI Environment Detection

#### Simulate GitHub Actions

```bash
docker run --rm \
  -e GITHUB_ACTIONS=true \
  -e GITHUB_REPOSITORY=test/repo \
  -e GITHUB_SHA=abc123 \
  -e GOVERNANCE_API_URL=http://host.docker.internal:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/test-data/sample-openapi.yaml \
  -e RULE_ID=test-rule-id \
  -v $(pwd):/workspace \
  governance-action
```

#### Simulate GitLab CI

```bash
docker run --rm \
  -e GITLAB_CI=true \
  -e CI_PROJECT_PATH=test/repo \
  -e CI_COMMIT_SHA=abc123 \
  -e GOVERNANCE_API_URL=http://host.docker.internal:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/test-data/sample-openapi.yaml \
  -e RULE_ID=test-rule-id \
  -v $(pwd):/workspace \
  governance-action
```

## Troubleshooting

### Common Issues

#### 1. Docker Connection Issues

**Problem**: Cannot connect to mock server from Docker container

**Solutions**:
```bash
# On macOS/Windows, use host.docker.internal
-e GOVERNANCE_API_URL=http://host.docker.internal:8080

# On Linux, try the Docker bridge IP
-e GOVERNANCE_API_URL=http://172.17.0.1:8080

# Or run with host network
docker run --rm --network host \
  -e GOVERNANCE_API_URL=http://localhost:8080 \
  # ... other environment variables
```

#### 2. File Path Issues

**Problem**: OAS file not found

**Solutions**:
```bash
# Verify the file exists
ls -la test-data/sample-openapi.yaml

# Use absolute path
-e OAS_FILE_PATH=/workspace/test-data/sample-openapi.yaml

# Check file permissions
chmod 644 test-data/sample-openapi.yaml
```

#### 3. Mock Server Not Responding

**Problem**: Mock server returns errors or doesn't start

**Solutions**:
```bash
# Check if port 8080 is available
lsof -i :8080

# Kill any process using the port
kill -9 $(lsof -t -i:8080)

# Start mock server with verbose logging
cd test-data
go run mock-server.go -verbose
```

#### 4. Permission Issues

**Problem**: Permission denied when reading files

**Solutions**:
```bash
# Check file permissions
ls -la test-data/

# Fix permissions
chmod 644 test-data/sample-openapi.yaml

# Run Docker with proper user mapping
docker run --rm -u $(id -u):$(id -g) \
  # ... other options
```

### Debug Mode

Enable verbose logging to get more detailed information:

```bash
# Set verbose mode
-e VERBOSE=true

# Check logs
docker logs <container_id>
```

### Network Debugging

Test network connectivity:

```bash
# Test from host
curl http://localhost:8080/health

# Test from Docker container
docker run --rm alpine/curl \
  curl http://host.docker.internal:8080/health
```

## Best Practices

1. **Always test locally** before committing changes
2. **Use the mock server** for consistent testing
3. **Test error scenarios** to ensure proper error handling
4. **Validate output format** matches expectations
5. **Test with different OAS files** to ensure compatibility
6. **Use version control** for test files and configurations
7. **Document test cases** for future reference

## Examples

### Complete Test Script

Create a test script for automated testing:

```bash
#!/bin/bash
# test-governance-action.sh

set -e

echo "Building Docker image..."
docker build -t governance-action .

echo "Starting mock server..."
cd test-data
go run mock-server.go &
MOCK_PID=$!
cd ..

sleep 2

echo "Testing governance action..."
docker run --rm \
  -e GOVERNANCE_API_URL=http://host.docker.internal:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/test-data/sample-openapi.yaml \
  -e RULE_ID=test-rule-id \
  -e VERBOSE=true \
  -v $(pwd):/workspace \
  governance-action

echo "Cleaning up..."
kill $MOCK_PID

echo "Test completed successfully!"
```

Make it executable and run:

```bash
chmod +x test-governance-action.sh
./test-governance-action.sh
```

This comprehensive testing approach ensures your governance action works correctly before deployment. 