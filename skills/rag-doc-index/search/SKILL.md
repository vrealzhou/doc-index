---
name: rag-doc-index-search
description: >
  Use this skill when you need to search indexed documentation using semantic search.
  This skill searches through already-indexed markdown documents and returns relevant
  sections based on natural language queries. Use it when the user asks questions about
  project documentation, needs to find specific information, or wants to understand
  concepts in the codebase. Requires documents to be indexed first using the index skill.
---

# RAG Document Search

Search indexed markdown documentation using semantic similarity.

## When to Use This Skill

Use this skill when you need to:
- Find relevant documentation sections in an indexed codebase
- Search for specific concepts, features, or patterns
- Understand project architecture or design decisions
- Retrieve context about APIs, configuration, or best practices
- Get answers to questions about the documentation

## Prerequisites

- **Indexed documents**: Documents must be indexed first using the `rag-doc-index-index` skill
- **TEI Server**: Must be running on localhost:8080
- **Configuration**: Must be configured with `doc-index config`

**Note**: If documents haven't been indexed yet, use the `rag-doc-index-index` skill first. See the main [rag-doc-index skill](../SKILL.md) for installation and setup instructions.

## Basic Usage

### Simple Search

```bash
# Basic search with natural language
doc-index search "authentication flow"

# Search for specific concepts
doc-index search "API rate limiting"

# Search for configuration details
doc-index search "database connection settings"
```

### Advanced Search Options

```bash
# Get more results
doc-index search "architecture" --top-k=10

# Filter by minimum score
doc-index search "API endpoints" --min-score=0.5

# JSON output for programmatic use
doc-index search "configuration" --json

# Combined options
doc-index search "microservices" --top-k=10 --min-score=0.4 --json
```

## Search Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--top-k=<n>` | integer | 5 | Maximum number of results to return |
| `--min-score=<n>` | float | 0.3 | Minimum similarity score (0-1) |
| `--json` | flag | false | Output results as JSON |

### Option Details

**`--top-k`**: Controls how many results to return
- Use higher values (10-15) for broad topics
- Use lower values (3-5) for specific queries
- Default: 5

**`--min-score`**: Filters results by relevance
- Range: 0.0 to 1.0 (higher = more relevant)
- Lower values (0.2-0.3) return more results, including less relevant ones
- Higher values (0.5-0.7) return only highly relevant results
- Default: 0.3

**`--json`**: Output format
- Returns machine-readable JSON
- Useful for scripting and programmatic access
- Default: human-readable format

## Understanding Search Results

### Human-Readable Output

```
Found 5 results:

1. [Score: 0.87] docs/architecture/overview.md#L15-42
   Section: System Architecture Overview
   This document describes the overall system architecture...

2. [Score: 0.82] docs/architecture/services.md#L8-25
   Section: Microservices Design
   Our system is composed of several microservices...
```

### JSON Output Format

```json
{
  "results": [
    {
      "document": "docs/architecture/overview.md",
      "section": "System Architecture Overview",
      "offset": 1024,
      "length": 512,
      "score": 0.87
    }
  ]
}
```

### Result Components

Each result contains:
- **Document**: Path to the source markdown file
- **Section**: Heading/title of the relevant section
- **Position**: Offset and length for precise content retrieval
- **Score**: Similarity score (0-1, higher is better)

## Common Use Cases

### 1. Finding Architecture Information

```bash
# Search for architecture documentation
doc-index search "system architecture" --top-k=5

# Search for specific architectural patterns
doc-index search "microservices design"
doc-index search "database architecture" --top-k=10

# Search for design decisions
doc-index search "why did we choose" --min-score=0.4
```

### 2. Understanding APIs

```bash
# Search for API documentation
doc-index search "API endpoints" --top-k=10

# Search for specific API features
doc-index search "REST API authentication"
doc-index search "rate limiting" --min-score=0.4

# Search for API examples
doc-index search "API usage examples"
```

### 3. Finding Configuration Details

```bash
# Search for configuration options
doc-index search "configuration" --top-k=5

# Search for environment variables
doc-index search "environment variables" --min-score=0.4

# Search for specific settings
doc-index search "database connection settings"
doc-index search "server configuration"
```

### 4. Debugging and Troubleshooting

```bash
# Search for error-related documentation
doc-index search "error handling" --top-k=5

# Search for specific error messages
doc-index search "connection timeout"
doc-index search "database connection failed" --min-score=0.3

# Search for debugging guides
doc-index search "troubleshooting guide"
```

### 5. Development Workflow

```bash
# Search for development guides
doc-index search "development setup" --top-k=5

# Search for testing documentation
doc-index search "testing guide" --min-score=0.5

# Search for contribution guidelines
doc-index search "contributing" --top-k=3

# Search for deployment information
doc-index search "deployment process"
```

### 6. Security and Best Practices

```bash
# Search for security documentation
doc-index search "security best practices" --top-k=8

# Search for authentication
doc-index search "authentication flow"
doc-index search "authorization" --min-score=0.4

# Search for performance
doc-index search "performance optimization"
```

## Tips for Effective Searches

### Use Natural Language

```bash
# Good - natural language queries
doc-index search "how to deploy the application"
doc-index search "error handling best practices"
doc-index search "how to configure the database"

# Avoid - too short or cryptic
doc-index search "deploy"
doc-index search "error"
doc-index search "config"
```

### Be Specific When Needed

```bash
# Broad search
doc-index search "API" --top-k=10

# More specific
doc-index search "REST API authentication"

# Very specific
doc-index search "JWT token validation for API requests"
```

### Adjust Score Thresholds

```bash
# Lower threshold if getting no results
doc-index search "rare topic" --min-score=0.2

# Default threshold for balanced results
doc-index search "common topic" --min-score=0.3

# Higher threshold for more precise results
doc-index search "important topic" --min-score=0.7
```

### Increase Result Count for Broad Topics

```bash
# Get more results for broad topics
doc-index search "architecture" --top-k=15

# Fewer results for specific queries
doc-index search "exact API endpoint" --top-k=3
```

### Try Different Query Formulations

```bash
# Try different ways to ask the same thing
doc-index search "user authentication"
doc-index search "login process"
doc-index search "user session management"
doc-index search "how users log in"
```

## Search Patterns

### Pattern 1: Broad to Narrow

Start with a broad search, then narrow down:

```bash
# Start broad
doc-index search "API" --top-k=10

# Narrow down
doc-index search "REST API authentication" --top-k=5

# Get specific
doc-index search "JWT token validation" --top-k=3
```

### Pattern 2: Multi-angle Search

Search from different angles to build complete context:

```bash
# Search for the main concept
doc-index search "database"

# Search for related concepts
doc-index search "database schema"
doc-index search "database migrations"
doc-index search "database performance"
```

### Pattern 3: Context Building

Build context by searching related topics:

```bash
doc-index search "authentication" --top-k=5
doc-index search "authorization" --top-k=5
doc-index search "user sessions" --top-k=5
doc-index search "security" --top-k=5
```

### Pattern 4: JSON for Scripting

Use JSON output for programmatic access:

```bash
# Get JSON and parse with jq
doc-index search "API" --json | jq '.results[0].document'

# Extract all document paths
doc-index search "error" --json | jq -r '.results[].document'

# Filter by score
doc-index search "config" --json | jq '.results[] | select(.score > 0.7)'
```

## Troubleshooting

### No Results Found

**Problem**: Search returns no results

**Solutions**:
1. **Lower the minimum score threshold**:
   ```bash
   doc-index search "query" --min-score=0.2
   ```

2. **Try different query terms**:
   ```bash
   # Instead of: doc-index search "auth"
   # Try: doc-index search "authentication"
   ```

3. **Check if documents are indexed**:
   ```bash
   doc-index status
   ```

4. **Verify configuration**:
   ```bash
   doc-index config --show
   ```

5. **Ensure TEI server is running**:
   ```bash
   curl http://localhost:8080/health
   ```

### Too Many Irrelevant Results

**Problem**: Getting too many irrelevant results

**Solutions**:
1. **Increase the minimum score**:
   ```bash
   doc-index search "query" --min-score=0.6
   ```

2. **Be more specific**:
   ```bash
   # Instead of: doc-index search "API"
   # Try: doc-index search "REST API authentication endpoints"
   ```

3. **Reduce the number of results**:
   ```bash
   doc-index search "query" --top-k=3
   ```

4. **Use more descriptive queries**:
   ```bash
   doc-index search "how to authenticate users with JWT tokens"
   ```

### Outdated Results

**Problem**: Search results seem outdated or missing recent changes

**Solutions**:
1. **Check document status**:
   ```bash
   doc-index status
   ```

2. **Reindex documents** (use the index skill):
   ```bash
   doc-index index --force
   ```

3. **Search again after reindexing**:
   ```bash
   doc-index search "updated content"
   ```

### TEI Connection Issues

**Problem**: Cannot connect to TEI server

**Solutions**:
1. **Verify TEI is running**:
   ```bash
   curl http://localhost:8080/health
   ```

2. **Check the endpoint URL**:
   ```bash
   doc-index config --show
   ```

3. **Restart TEI server** if needed

4. **Check firewall settings**

### Poor Search Quality

**Problem**: Search results are not relevant

**Solutions**:
1. **Use more descriptive queries**:
   ```bash
   # Instead of: doc-index search "error"
   # Try: doc-index search "how to handle database connection errors"
   ```

2. **Adjust min-score threshold**:
   ```bash
   # Try different thresholds
   doc-index search "query" --min-score=0.4
   ```

3. **Increase top-k for more options**:
   ```bash
   doc-index search "query" --top-k=10
   ```

4. **Try different query formulations**:
   ```bash
   doc-index search "user login"
   doc-index search "authentication process"
   doc-index search "user session management"
   ```

## Quick Reference

| Goal | Command |
|------|---------|
| Basic search | `doc-index search "query"` |
| More results | `doc-index search "query" --top-k=10` |
| Filter by score | `doc-index search "query" --min-score=0.5` |
| JSON output | `doc-index search "query" --json` |
| Combined options | `doc-index search "query" --top-k=10 --min-score=0.5 --json` |

## Related Skills

- **rag-doc-index-index**: Use this skill to index documents before searching
- **rag-doc-index** (parent): Overview and installation instructions for the RAG document indexing system

## Notes

- Documents must be indexed before searching (use the index skill)
- TEI server must be running for search operations
- Use natural language queries for best results
- Lower `--min-score` (default: 0.3) if getting no results
- Higher `--min-score` for more precise results
- The search uses semantic similarity, not keyword matching
- Results include position information for precise content retrieval