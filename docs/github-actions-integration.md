# GitHub Actions Integration

This document provides detailed instructions for integrating the Governance Action with GitHub Actions workflows.

## Overview

The Governance Action is fully compatible with GitHub Actions and can be used to enforce API governance standards in your GitHub repositories. It supports:

- Automatic GitHub Actions environment detection
- GitHub-specific context extraction
- Output variable generation for downstream steps
- Integration with GitHub's security features

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
        uses: docker://ghcr.io/tyktechnologies/governance-action:latest
        env:
          GOVERNANCE_API_URL: ${{ secrets.GOVERNANCE_API_URL }}
          GOVERNANCE_API_TOKEN: ${{ secrets.GOVERNANCE_API_TOKEN }}
          OAS_FILE_PATH: ./api/openapi.yaml
          RULE_ID: ${{ secrets.GOVERNANCE_RULE_ID }}
          VERBOSE: true
```

### 2. Configure Repository Secrets

In your GitHub repository:

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Add the following repository secrets:
   - `GOVERNANCE_API_URL`: Your governance service URL
   - `GOVERNANCE_API_TOKEN`: Your governance service authentication token
   - `GOVERNANCE_RULE_ID`: The rule ID to evaluate against

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
        uses: docker://ghcr.io/tyktechnologies/governance-action:latest
        env:
          GOVERNANCE_API_URL: ${{ secrets.GOVERNANCE_API_URL }}
          GOVERNANCE_API_TOKEN: ${{ secrets.GOVERNANCE_API_TOKEN }}
          OAS_FILE_PATH: ./api/openapi.yaml
          RULE_ID: ${{ secrets.GOVERNANCE_RULE_ID }}
          VERBOSE: true

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
        uses: docker://ghcr.io/tyktechnologies/governance-action:latest
        env:
          GOVERNANCE_API_URL: ${{ secrets.GOVERNANCE_API_URL }}
          GOVERNANCE_API_TOKEN: ${{ secrets.GOVERNANCE_API_TOKEN }}
          OAS_FILE_PATH: ./api/openapi.yaml
          RULE_ID: ${{ secrets.GOVERNANCE_RULE_ID }}
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
        uses: docker://ghcr.io/tyktechnologies/governance-action:latest
        env:
          GOVERNANCE_API_URL: ${{ secrets.GOVERNANCE_API_URL }}
          GOVERNANCE_API_TOKEN: ${{ secrets.GOVERNANCE_API_TOKEN }}
          OAS_FILE_PATH: ./api/openapi.yaml
          RULE_ID: ${{ secrets.API_GOVERNANCE_RULE_ID }}

  governance-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Check Documentation Governance
        uses: docker://ghcr.io/tyktechnologies/governance-action:latest
        env:
          GOVERNANCE_API_URL: ${{ secrets.GOVERNANCE_API_URL }}
          GOVERNANCE_API_TOKEN: ${{ secrets.GOVERNANCE_API_TOKEN }}
          OAS_FILE_PATH: ./docs/api.yaml
          RULE_ID: ${{ secrets.DOCS_GOVERNANCE_RULE_ID }}
```

## Environment Variables

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `GOVERNANCE_API_URL` | Base URL of the governance service | `https://governance.example.com` |
| `GOVERNANCE_API_TOKEN` | Authentication token for the governance API | `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...` |
| `OAS_FILE_PATH` | Path to the OpenAPI Specification file | `./api/openapi.yaml` |
| `RULE_ID` | ID of the governance rule to evaluate against | `6853d42c7493327ea805be8a` |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VERBOSE` | Enable verbose logging | `false` |

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
  uses: docker://ghcr.io/tyktechnologies/governance-action:latest
  env:
    GOVERNANCE_API_URL: ${{ secrets.GOVERNANCE_API_URL }}
    GOVERNANCE_API_TOKEN: ${{ secrets.GOVERNANCE_API_TOKEN }}
    OAS_FILE_PATH: ./api/openapi.yaml
    RULE_ID: ${{ secrets.GOVERNANCE_RULE_ID }}

- name: Check Results
  run: |
    echo "Found ${{ steps.governance.outputs.error_count }} errors"
    echo "Found ${{ steps.governance.outputs.warning_count }} warnings"
    echo "Total issues: ${{ steps.governance.outputs.total_issues }}"
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
env:
  GOVERNANCE_API_URL: ${{ secrets.GOVERNANCE_API_URL }}
  GOVERNANCE_API_TOKEN: ${{ secrets.GOVERNANCE_API_TOKEN }}
  RULE_ID: ${{ secrets.GOVERNANCE_RULE_ID }}
```

### Environment-Specific Secrets

For different environments, use environment-specific secrets:

```yaml
jobs:
  governance:
    environment: production
    steps:
      - name: Run Governance Check
        uses: docker://ghcr.io/tyktechnologies/governance-action:latest
        env:
          GOVERNANCE_API_URL: ${{ secrets.PROD_GOVERNANCE_API_URL }}
          GOVERNANCE_API_TOKEN: ${{ secrets.PROD_GOVERNANCE_API_TOKEN }}
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
   curl -H "Authorization: Bearer $GOVERNANCE_API_TOKEN" $GOVERNANCE_API_URL/health
   ```

3. **File Not Found**
   ```bash
   # Verify the OAS file path is correct
   find . -name "*.yaml" -o -name "*.yml" -o -name "*.json"
   ```

### Debug Mode

Enable verbose logging to troubleshoot issues:

```yaml
env:
  VERBOSE: true
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

### Advanced Workflows

For more complex workflows, check out the examples in the [GitHub Actions documentation](https://docs.github.com/en/actions). 