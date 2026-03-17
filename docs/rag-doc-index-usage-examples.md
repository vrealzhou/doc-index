# Verify TEI is running
curl http://localhost:8080/health

# Index documents (if not already done)
doc-index index
```

## Example 1: Finding Architecture Information

**Scenario**: You need to understand the overall architecture of a project.

```bash
# Search for architecture documentation
doc-index search "system architecture" --top-k=5

# Search for specific components
doc-index search "microservices" --min-score=0.4
doc-index search "database design" --top-k=10
```

**Expected Output**:
```
Found 5 results:

1. [Score: 0.87] docs/architecture/overview.md#L15-42
   Section: System Architecture Overview
   This document describes the overall system architecture...

2. [Score: 0.82] docs/architecture/services.md#L8-25
   Section: Microservices Design
   Our system is composed of several microservices...
```

## Example 2: Understanding Authentication Flow

**Scenario**: You need to understand how authentication works in the codebase.

```bash
# Search for authentication documentation
doc-index search "authentication flow" --top-k=5

# Get more specific results
doc-index search "user login process" --min-score=0.5
doc-index search "JWT token validation"
```

## Example 3: API Documentation Search

**Scenario**: You need to find API endpoint documentation.

```bash
# Search for API endpoints
doc-index search "API endpoints" --top-k=10

# Search for specific API features
doc-index search "rate limiting" --min-score=0.4
doc-index search "REST API authentication"
```

## Example 4: Database Schema Discovery

**Scenario**: You need to understand the database structure.

```bash
# Search for database schema information
doc-index search "database schema" --top-k=5

# Search for specific tables or models
doc-index search "user table structure"
doc-index search "database migrations" --min-score=0.5
```

## Example 5: JSON Output for Programmatic Use

**Scenario**: You want to use the results in a script or another tool.

```bash
# Get JSON output
doc-index search "configuration options" --json

# Parse with jq
doc-index search "deployment guide" --json | jq '.results[0].document'
```

**JSON Output Format**:
```json
{
  "results": [
    {
      "document": "docs/deployment/guide.md",
      "section": "Deployment Guide",
      "position": {
        "offset": 245,
        "length": 1523
      },
      "score": 0.91
    }
  ]
}
```

## Example 6: Finding Configuration Details

**Scenario**: You need to find configuration information.

```bash
# Search for configuration documentation
doc-index search "configuration" --top-k=5

# Search for environment variables
doc-index search "environment variables" --min-score=0.4

# Search for specific config options
doc-index search "database connection settings"
```

## Example 7: Troubleshooting Documentation

**Scenario**: You're encountering an error and need to find troubleshooting docs.

```bash
# Search for error-related documentation
doc-index search "error handling" --top-k=5

# Search for specific error messages
doc-index search "connection timeout"
doc-index search "database connection failed" --min-score=0.3
```

## Example 8: Development Workflow

**Scenario**: You want to understand the development workflow.

```bash
# Search for development guides
doc-index search "development setup" --top-k=5

# Search for testing documentation
doc-index search "testing guide" --min-score=0.5

# Search for contribution guidelines
doc-index search "contributing" --top-k=3
```

## Example 9: Multi-concept Search

**Scenario**: You need documentation covering multiple related concepts.

```bash
# Search for security-related documentation
doc-index search "security best practices" --top-k=8

# Search for performance optimization
doc-index search "performance optimization" --min-score=0.4

# Search for monitoring and logging
doc-index search "monitoring and logging" --top-k=6
```

## Example 10: Reindexing After Updates

**Scenario**: Documentation has been updated and you need to reindex.

```bash
# Check current status
doc-index status

# Force reindex to capture changes
doc-index index --force

# Verify the update
doc-index search "new feature" --top-k=5
```

## Tips for Effective Searches

### Use Natural Language
```bash
# Good - natural language queries
doc-index search "how to deploy the application"
doc-index search "error handling best practices"

# Avoid - too short or cryptic
doc-index search "deploy"
doc-index search "error"
```

### Adjust Score Thresholds
```bash
# Lower threshold if getting no results
doc-index search "rare topic" --min-score=0.2

# Higher threshold for more precise results
doc-index search "common topic" --min-score=0.7
```

### Increase Result Count
```bash
# Get more results for broad topics
doc-index search "architecture" --top-k=15

# Fewer results for specific queries
doc-index search "exact API endpoint" --top-k=3
```

### Combine with Other Tools
```bash
# Use with grep for file-level search
doc-index search "authentication" --json | jq -r '.results[].document' | xargs grep "JWT"

# Use with file reading
doc-index search "configuration" --top-k=1 | head -1 | cut -d: -f1
```

## Common Patterns

### Pattern 1: Broad to Narrow
```bash
# Start broad, then narrow down
doc-index search "API" --top-k=10
doc-index search "REST API authentication" --top-k=5
doc-index search "JWT token validation" --top-k=3
```

### Pattern 2: Multi-angle Search
```bash
# Search from different angles
doc-index search "user authentication"
doc-index search "login process"
doc-index search "user session management"
```

### Pattern 3: Context Building
```bash
# Build context by searching related topics
doc-index search "database schema" --top-k=5
doc-index search "database migrations" --top-k=5
doc-index search "database performance" --top-k=5
```

## Troubleshooting Examples

### No Results Found
```bash
# If getting no results, try:
# 1. Lower the minimum score
doc-index search "query" --min-score=0.2

# 2. Use different keywords
doc-index search "alternative terms"

# 3. Check if documents are indexed
doc-index status
```

### Too Many Irrelevant Results
```bash
# If getting too many irrelevant results:
# 1. Increase the minimum score
doc-index search "query" --min-score=0.6

# 2. Be more specific
doc-index search "more specific query terms"

# 3. Reduce the number of results
doc-index search "query" --top-k=3
```

### Outdated Results
```bash
# If results seem outdated:
# 1. Check document status
doc-index status

# 2. Force reindex
doc-index index --force

# 3. Search again
doc-index search "query" --top-k=5