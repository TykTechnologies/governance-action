# GitHub Actions Integration

This document provides detailed instructions for integrating the Governance Action with GitHub Actions workflows.

## Overview

The Governance Action is fully compatible with GitHub Actions and can be used to enforce API governance standards in your GitHub repositories. It supports:

- Automatic GitHub Actions environment detection
- GitHub-specific context extraction
- Output variable generation for downstream steps
- Integration with GitHub's security features
- **Mock mode for testing and development**

## Quick Start

### 1. Add the Workflow Configuration

Create a `.github/workflows/governance.yml` file in your repository:

```yaml
name: Governance Check
on: [push, pull_request]

jobs:
  governance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Governance Check
        uses: tyktechnologies/governance-action@latest
        with:
          governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
          governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
          rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
          api_path: ./api/openapi.yaml
```

### 2. Configure Repository Secrets

In your GitHub repository:

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Add the following repository secrets:
   - `GOVERNANCE_SERVICE_URL`: Your governance service URL
   - `GOVERNANCE_SERVICE_TOKEN`: Your governance service authentication token
   - `GOVERNANCE_RULE_ID`: The rule ID to evaluate against

## Mock Mode for Testing

The action supports a **mock mode** that bypasses the governance service API call and returns predefined results. This is perfect for:

- **Development and testing** without a real governance service
- **CI/CD pipeline testing** to verify different scenarios
- **Documentation examples** and demonstrations

### Mock Mode Usage

```yaml
- name: Test Governance Check (Mock Mode)
  uses: tyktechnologies/governance-action@latest
  with:
    rule_id: test-rule-id
    api_path: ./api/openapi.yaml
    mocked: success  # Options: success, fail, warning
```

### Available Mock Scenarios

| Mock Value | Description | Exit Code | Use Case |
|------------|-------------|-----------|----------|
| `success` | No governance issues found | 0 | Test successful scenarios |
| `warning` | 2 warnings, no errors | 0 | Test warning-only scenarios |
| `fail` | 2 errors + 1 warning | 1 | Test failure scenarios |

### Complete Testing Workflow

```yaml
name: Governance Testing
on: [push, pull_request]

jobs:
  test-governance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Test Success Scenario
        uses: tyktechnologies/governance-action@latest
        with:
          rule_id: test-rule-id
          api_path: ./api/openapi.yaml
          mocked: success
          
      - name: Test Warning Scenario
        uses: tyktechnologies/governance-action@latest
        with:
          rule_id: test-rule-id
          api_path: ./api/openapi.yaml
          mocked: warning
          
      - name: Test Fail Scenario
        uses: tyktechnologies/governance-action@latest
        continue-on-error: true  # Don't fail the workflow for this test
        with:
          rule_id: test-rule-id
          api_path: ./api/openapi.yaml
          mocked: fail
```

### Conditional Mock Usage

```yaml
- name: Run Governance Check
  uses: tyktechnologies/governance-action@latest
  with:
    governance_service: ${{ github.event_name == 'push' && secrets.GOVERNANCE_SERVICE_URL || '' }}
    governance_auth: ${{ github.event_name == 'push' && secrets.GOVERNANCE_SERVICE_TOKEN || '' }}
    rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
    api_path: ./api/openapi.yaml
    mocked: ${{ github.event_name == 'pull_request' && 'warning' || '' }}
```

## Advanced Configuration

### Multi-Step Workflow

```yaml
name: CI/CD Pipeline
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run tests
        run: |
          echo "Running tests..."
          # Your test commands here

  governance:
    runs-on: ubuntu-latest
    needs: test
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Governance Check
        id: governance
        uses: tyktechnologies/governance-action@latest
        with:
          governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
          governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
          rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
          api_path: ./api/openapi.yaml

  deploy:
    runs-on: ubuntu-latest
    needs: [test, governance]
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - name: Deploy
        run: |
          echo "Found ${{ steps.governance.outputs.error_count }} errors"
          echo "Found ${{ steps.governance.outputs.warning_count }} warnings"
          echo "Deploying application..."
```

### Conditional Execution

```yaml
name: Governance Check
on: [push, pull_request]

jobs:
  governance:
    runs-on: ubuntu-latest
    if: |
      github.event_name == 'pull_request' ||
      github.ref == 'refs/heads/main' ||
      github.ref == 'refs/heads/develop'
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Governance Check
        uses: tyktechnologies/governance-action@latest
        with:
          governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
          governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
          rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
          api_path: ./api/openapi.yaml
```

### Multiple API Files

```yaml
name: Governance Check
on: [push, pull_request]

jobs:
  governance-api:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Check API Governance
        uses: tyktechnologies/governance-action@latest
        with:
          governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
          governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
          rule_id: ${{ secrets.API_GOVERNANCE_RULE_ID }}
          api_path: ./api/openapi.yaml

  governance-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Check Documentation Governance
        uses: tyktechnologies/governance-action@latest
        with:
          governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
          governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
          rule_id: ${{ secrets.DOCS_GOVERNANCE_RULE_ID }}
          api_path: ./docs/api.yaml
```

## Input Parameters

### Required Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `governance_service` | Base URL of the governance service | `https://governance.example.com` |
| `governance_auth` | Authentication token for the governance API | `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...` |
| `rule_id` | ID of the governance rule to evaluate against | `6853d42c7493327ea805be8a` |
| `api_path` | Path to the OpenAPI Specification file | `./api/openapi.yaml` |

### Optional Parameters

| Parameter | Description | Default | Notes |
|-----------|-------------|---------|-------|
| `mocked` | Mock mode for testing (`success`, `fail`, `warning`) | - | When set, bypasses API call |

### Environment Variable Fallbacks

The action also supports environment variables:

| Environment Variable | Maps to Input Parameter |
|---------------------|-------------------------|
| `GOVERNANCE_SERVICE` | `governance_service` |
| `GOVERNANCE_AUTH` | `governance_auth` |
| `RULE_ID` | `rule_id` |
| `API_PATH` | `api_path` |
| `MOCKED` | `mocked` |

### GitHub-Specific Variables

The action automatically detects and uses these GitHub Actions variables:

| Variable | Description |
|----------|-------------|
| `GITHUB_REPOSITORY` | Repository name (e.g., `owner/repo`) |
| `GITHUB_SHA` | Commit SHA |
| `GITHUB_REF_NAME` | Branch or tag name |
| `GITHUB_ACTOR` | User who triggered the workflow |
| `GITHUB_WORKFLOW` | Workflow name |
| `GITHUB_RUN_ID` | Run ID |

## Output Variables

The action provides the following outputs that can be used in subsequent steps:

### Using Outputs

```yaml
- name: Run Governance Check
  id: governance
  uses: tyktechnologies/governance-action@latest
  with:
    governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
    governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
    rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
    api_path: ./api/openapi.yaml

- name: Check Results
  run: |
    echo "Found ${{ steps.governance.outputs.error_count }} errors"
    echo "Found ${{ steps.governance.outputs.warning_count }} warnings"
    echo "Total issues: ${{ steps.governance.outputs.total_issues }}"
```

### Using Outputs with Mock Mode

```yaml
- name: Test Governance Check
  id: governance-test
  uses: tyktechnologies/governance-action@latest
  with:
    rule_id: test-rule-id
    api_path: ./api/openapi.yaml
    mocked: fail

- name: Check Mock Results
  run: |
    echo "Mock test found ${{ steps.governance-test.outputs.error_count }} errors"
    echo "Mock test found ${{ steps.governance-test.outputs.warning_count }} warnings"
    echo "Mock test total issues: ${{ steps.governance-test.outputs.total_issues }}"
```

### Available Outputs

| Output | Description |
|--------|-------------|
| `error_count` | Number of governance errors found |
| `warning_count` | Number of governance warnings found |
| `total_issues` | Total number of governance issues found |

## Security Best Practices

### Using GitHub Secrets

Always use GitHub secrets for sensitive information:

```yaml
with:
  governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
  governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
  rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
```

### Environment-Specific Secrets

For different environments, use environment-specific secrets:

```yaml
jobs:
  governance:
    environment: production
    steps:
      - name: Run Governance Check
        uses: tyktechnologies/governance-action@latest
        with:
          governance_service: ${{ secrets.PROD_GOVERNANCE_SERVICE_URL }}
          governance_auth: ${{ secrets.PROD_GOVERNANCE_SERVICE_TOKEN }}
          rule_id: ${{ secrets.PROD_GOVERNANCE_RULE_ID }}
          api_path: ./api/openapi.yaml
```

### Branch Protection

Configure branch protection rules to require governance checks:

1. Go to **Settings** → **Branches**
2. Add a rule for your main branch
3. Check "Require status checks to pass before merging"
4. Add your governance workflow as a required status check

## Troubleshooting

### Common Issues

1. **Permission Denied**
   ```bash
   # Ensure the action has access to your files
   ls -la ./api/openapi.yaml
   ```

2. **Network Connectivity**
   ```bash
   # Test connectivity to governance service
   curl -H "Authorization: Bearer $GOVERNANCE_SERVICE_TOKEN" $GOVERNANCE_SERVICE_URL/health
   ```

3. **File Not Found**
   ```bash
   # Verify the OAS file path is correct
   find . -name "*.yaml" -o -name "*.yml" -o -name "*.json"
   ```

4. **Mock Mode Issues**
   ```yaml
   # Ensure mocked parameter is one of: success, fail, warning
   with:
     rule_id: test-rule-id
     api_path: ./api/openapi.yaml
     mocked: success  # Valid values only
   ```

### Debug Mode

Enable verbose logging to troubleshoot issues:

```yaml
# For normal mode, use environment variable
env:
  VERBOSE: true

# For mocked mode, use mocked parameter
with:
  rule_id: test-rule-id
  api_path: ./api/openapi.yaml
  mocked: warning  # This will show detailed output
```

### Local Testing

Test your GitHub Actions workflow locally using `act`:

```bash
# Install act
brew install act  # macOS
# or
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash  # Linux

# Run the workflow locally
act push
```

## Examples

### Basic Workflow

See the [`.github/workflows/governance.yml`](../.github/workflows/governance.yml) file in this repository for a complete example configuration.

### Complete Workflow Examples

#### Production Workflow

```yaml
name: Production Governance Check
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  governance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Governance Check
        id: governance
        uses: tyktechnologies/governance-action@latest
        with:
          governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
          governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
          rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
          api_path: ./api/openapi.yaml

      - name: Deploy on Success
        if: steps.governance.outputs.error_count == '0'
        run: |
          echo "No governance errors found. Proceeding with deployment..."
          # Your deployment commands here
```

#### Development Workflow with Testing

```yaml
name: Development Governance Testing
on:
  push:
    branches: [ develop, feature/* ]
  pull_request:
    branches: [ develop ]

jobs:
  test-governance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Test Success Scenario
        uses: tyktechnologies/governance-action@latest
        with:
          rule_id: test-rule-id
          api_path: ./api/openapi.yaml
          mocked: success
          
      - name: Test Warning Scenario
        uses: tyktechnologies/governance-action@latest
        with:
          rule_id: test-rule-id
          api_path: ./api/openapi.yaml
          mocked: warning
          
      - name: Test Fail Scenario
        uses: tyktechnologies/governance-action@latest
        continue-on-error: true
        with:
          rule_id: test-rule-id
          api_path: ./api/openapi.yaml
          mocked: fail

  real-governance:
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Real Governance Check
        uses: tyktechnologies/governance-action@latest
        with:
          governance_service: ${{ secrets.GOVERNANCE_SERVICE_URL }}
          governance_auth: ${{ secrets.GOVERNANCE_SERVICE_TOKEN }}
          rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
          api_path: ./api/openapi.yaml
```

#### Conditional Workflow

```yaml
name: Smart Governance Check
on: [push, pull_request]

jobs:
  governance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Governance Check
        uses: tyktechnologies/governance-action@latest
        with:
          governance_service: ${{ github.event_name == 'push' && secrets.GOVERNANCE_SERVICE_URL || '' }}
          governance_auth: ${{ github.event_name == 'push' && secrets.GOVERNANCE_SERVICE_TOKEN || '' }}
          rule_id: ${{ secrets.GOVERNANCE_RULE_ID }}
          api_path: ./api/openapi.yaml
          mocked: ${{ github.event_name == 'pull_request' && 'warning' || '' }}
```

### Advanced Workflows

For more complex workflows, check out the examples in the [GitHub Actions documentation](https://docs.github.com/en/actions). 