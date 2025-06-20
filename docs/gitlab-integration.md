# GitLab CI Integration

This document provides detailed instructions for integrating the Governance Action with GitLab CI/CD pipelines.

## Overview

The Governance Action is fully compatible with GitLab CI and can be used to enforce API governance standards in your GitLab projects. It supports:

- Automatic GitLab CI environment detection
- GitLab-specific context extraction
- Output variable generation for downstream jobs
- Artifact generation for test reporting

## Quick Start

### 1. Add the Pipeline Configuration

Copy the following configuration to your `.gitlab-ci.yml` file:

```yaml
stages:
  - governance

governance-check:
  stage: governance
  image: ghcr.io/tyktechnologies/governance-action:latest
  variables:
    GOVERNANCE_API_URL: $GOVERNANCE_API_URL
    GOVERNANCE_API_TOKEN: $GOVERNANCE_API_TOKEN
    OAS_FILE_PATH: ./api/openapi.yaml
    RULE_ID: $GOVERNANCE_RULE_ID
    VERBOSE: "true"
  script:
    - /app/governance-action
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

### 2. Configure Environment Variables

In your GitLab project:

1. Go to **Settings** → **CI/CD** → **Variables**
2. Add the following variables:
   - `GOVERNANCE_API_URL`: Your governance service URL
   - `GOVERNANCE_API_TOKEN`: Your governance service authentication token
   - `GOVERNANCE_RULE_ID`: The rule ID to evaluate against

## Advanced Configuration

### Multi-Stage Pipeline

```yaml
stages:
  - test
  - governance
  - deploy

test:
  stage: test
  script:
    - echo "Running tests..."

governance-check:
  stage: governance
  image: ghcr.io/tyktechnologies/governance-action:latest
  variables:
    GOVERNANCE_API_URL: $GOVERNANCE_API_URL
    GOVERNANCE_API_TOKEN: $GOVERNANCE_API_TOKEN
    OAS_FILE_PATH: ./api/openapi.yaml
    RULE_ID: $GOVERNANCE_RULE_ID
    VERBOSE: "true"
    GITLAB_OUTPUT_FILE: governance_output.env
  script:
    - /app/governance-action
  artifacts:
    paths:
      - governance_output.env
    expire_in: 1 week
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH

deploy:
  stage: deploy
  script:
    - source governance_output.env
    - echo "Found $error_count errors and $warning_count warnings"
    - if [ "$error_count" -gt 0 ]; then echo "Cannot deploy due to governance errors"; exit 1; fi
    - echo "Deploying application..."
  dependencies:
    - governance-check
  rules:
    - if: $CI_COMMIT_BRANCH == $CI_DEFAULT_BRANCH
```

### Conditional Execution

```yaml
governance-check:
  stage: governance
  image: ghcr.io/tyktechnologies/governance-action:latest
  variables:
    GOVERNANCE_API_URL: $GOVERNANCE_API_URL
    GOVERNANCE_API_TOKEN: $GOVERNANCE_API_TOKEN
    OAS_FILE_PATH: ./api/openapi.yaml
    RULE_ID: $GOVERNANCE_RULE_ID
  script:
    - /app/governance-action
  rules:
    # Only run on merge requests
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
    # Only run on main branch
    - if: $CI_COMMIT_BRANCH == "main"
    # Skip if governance is disabled
    - if: $GOVERNANCE_DISABLED == "true"
      when: never
    # Run on all other branches
    - when: on_success
```

### Parallel Jobs

```yaml
governance-check-api:
  stage: governance
  image: ghcr.io/tyktechnologies/governance-action:latest
  variables:
    GOVERNANCE_API_URL: $GOVERNANCE_API_URL
    GOVERNANCE_API_TOKEN: $GOVERNANCE_API_TOKEN
    OAS_FILE_PATH: ./api/openapi.yaml
    RULE_ID: $API_GOVERNANCE_RULE_ID
  script:
    - /app/governance-action

governance-check-docs:
  stage: governance
  image: ghcr.io/tyktechnologies/governance-action:latest
  variables:
    GOVERNANCE_API_URL: $GOVERNANCE_API_URL
    GOVERNANCE_API_TOKEN: $GOVERNANCE_API_TOKEN
    OAS_FILE_PATH: ./docs/api.yaml
    RULE_ID: $DOCS_GOVERNANCE_RULE_ID
  script:
    - /app/governance-action
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
| `GITLAB_OUTPUT_FILE` | Path for output variables file | `governance_output.env` |

### GitLab-Specific Variables

The action automatically detects and uses these GitLab CI variables:

| Variable | Description |
|----------|-------------|
| `CI_PROJECT_PATH` | Full project path (e.g., `group/project`) |
| `CI_COMMIT_SHA` | Commit SHA |
| `CI_COMMIT_BRANCH` | Branch name |
| `GITLAB_USER_NAME` | User who triggered the pipeline |
| `CI_PIPELINE_ID` | Pipeline ID |
| `CI_JOB_ID` | Job ID |

## Output Variables

The action generates the following output variables:

### Environment Variables

These are available in the current job:

```bash
echo "Found $error_count errors"
echo "Found $warning_count warnings"
echo "Total issues: $total_issues"
```

### Output File

If `GITLAB_OUTPUT_FILE` is set, variables are exported to a file:

```bash
# governance_output.env
export error_count=2
export warning_count=1
export total_issues=3
```

### Using Outputs in Downstream Jobs

```yaml
deploy:
  stage: deploy
  script:
    - source governance_output.env
    - if [ "$error_count" -gt 0 ]; then
        echo "Cannot deploy: $error_count governance errors found"
        exit 1
      fi
    - echo "Deploying with $warning_count warnings"
  dependencies:
    - governance-check
```

## Troubleshooting

### Common Issues

1. **Permission Denied**
   ```bash
   # Ensure the container has read access to your files
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
variables:
  VERBOSE: "true"
```

### Local Testing

Test your GitLab CI configuration locally using Docker:

```bash
docker run --rm \
  -e GITLAB_CI=true \
  -e GOVERNANCE_API_URL=http://localhost:8080 \
  -e GOVERNANCE_API_TOKEN=mock-token \
  -e OAS_FILE_PATH=/workspace/api/openapi.yaml \
  -e RULE_ID=test-rule-id \
  -v $(pwd):/workspace \
  ghcr.io/tyktechnologies/governance-action:latest
```

## Best Practices

1. **Security**: Use GitLab's protected variables for sensitive information
2. **Performance**: Run governance checks early in the pipeline
3. **Flexibility**: Use conditional rules to skip governance when appropriate
4. **Monitoring**: Set up alerts for governance failures
5. **Documentation**: Keep governance rules and requirements documented

## Examples

See the [`.gitlab-ci.yml`](../.gitlab-ci.yml) file in this repository for a complete example configuration. 