# API Reference

Complete reference for all `doc-index` CLI commands.

## Commands Overview

| Command | Description |
|---------|-------------|
| `config` | Configure TEI endpoint and documents path |
| `index` | Index markdown documents |
| `search` | Search indexed documents |
| `status` | Check indexing status |

---

## config

Configure the skill with TEI endpoint and documents path.

### Usage

```bash
doc-index config [options]
```

### Options

| Option | Type | Description |
|--------|------|-------------|
| `--tei=<url>` | string | TEI server endpoint URL |
| `--docs=<path>` | string | Root path to markdown documents |
| `--show` | flag | Display current configuration |

### Examples

```bash
# Set TEI endpoint and docs path
doc-index config --tei=http://localhost:8080 --docs=/path/to/docs

# Show current configuration
doc-index config --show
```

### Output

When using `--show`:
```
TEI Endpoint: http://localhost:8080
Documents Path: /path/to/docs
```

---

## index

Index all markdown documents in the configured path.

### Usage

```bash
doc-index index [options]
```

### Options

| Option | Type | Description |
|--------|------|-------------|
| `--force` | flag | Force reindex all documents |

### Behavior

- Scans all `.md` files in the configured `--docs` path
- Chunks documents by headers and sections
- Generates embeddings via TEI server
- Stores embeddings in JSONL format in sibling `embeddings/` directory
- Skips documents that haven't changed (hash-based detection)

### Examples

```bash
# Index new/modified documents
doc-index index

# Force reindex all documents
doc-index index --force
```

### Output

```
Indexing documents...
Found 42 markdown files
Indexed 38 new/modified documents
Skipped 4 unchanged documents
Total chunks: 156
Embeddings saved to: /path/to/embeddings/
```

---

## search

Search indexed documents using semantic similarity.

### Usage

```bash
doc-index search <query> [options]
```

### Arguments

| Argument | Type | Description |
|----------|------|-------------|
| `query` | string | Natural language search query (required) |

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `--top-k=<n>` | integer | 5 | Maximum number of results to return |
| `--min-score=<n>` | float | 0.3 | Minimum similarity score (0-1) |
| `--json` | flag | false | Output results as JSON |

### Output Format

**Default (human-readable):**
```
Result 1 (score: 0.85)
  Document: docs/api/authentication.md
  Section: JWT Token Validation
  Position: offset=1024, length=512

Result 2 (score: 0.78)
  Document: docs/guides/security.md
  Section: Token Expiration
  Position: offset=2048, length=256
...
```

**JSON format (`--json`):**
```json
{
  "results": [
    {
      "document": "docs/api/authentication.md",
      "section": "JWT Token Validation",
      "offset": 1024,
      "length": 512,
      "score": 0.85
    },
    {
      "document": "docs/guides/security.md",
      "section": "Token Expiration",
      "offset": 2048,
      "length": 256,
      "score": 0.78
    }
  ]
}
```

### Examples

```bash
# Basic search
doc-index search "authentication flow"

# Get more results
doc-index search "API rate limiting" --top-k=10

# Filter by minimum score
doc-index search "database schema" --min-score=0.5

# JSON output for programmatic use
doc-index search "error handling" --json

# Combined options
doc-index search "API endpoints" --top-k=10 --min-score=0.5 --json
```

---

## status

Check the indexing status and statistics.

### Usage

```bash
doc-index status
```

### Output

```
Index Status:
  Documents indexed: 42
  Total chunks: 156
  Last indexed: 2024-01-15 14:32:00
  Embeddings location: /path/to/embeddings/
  
TEI Server:
  Status: connected
  Endpoint: http://localhost:8080
  Model: BAAI/bge-small-en-v1.5
```

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error |
| 3 | TEI connection error |
| 4 | Index not found |

---

## Error Messages

### Configuration Errors

```
Error: Configuration file not found
Solution: Run 'doc-index config --tei=<url> --docs=<path>'
```

```
Error: Invalid TEI endpoint
Solution: Ensure TEI server is running and URL is correct
```

### Index Errors

```
Error: No documents indexed
Solution: Run 'doc-index index' to index documents
```

```
Error: Embeddings directory not found
Solution: Check if index operation completed successfully
```

### Search Errors

```
Error: Query too short
Solution: Provide a more descriptive search query
```

```
Error: No results found
Solution: Try lowering --min-score or use different query terms
```

---

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DOC_INDEX_CONFIG` | Custom config file path | `./rag-doc-index.config.json` |
| `TEI_TIMEOUT` | TEI server timeout (seconds) | `30` |
| `MAX_CHUNK_SIZE` | Maximum chunk size (characters) | `512` |

---

## Configuration File Schema

**Location:** `rag-doc-index.config.json`

```json
{
  "tei_endpoint": "http://localhost:8080",
  "docs_path": "/path/to/project/docs",
  "model_id": "BAAI/bge-small-en-v1.5",
  "chunk_size": 512,
  "chunk_overlap": 50
}
```

---

## Storage Format

### Embeddings Storage

**Location:** `<docs_path>/../embeddings/*.jsonl`

**Format:** JSONL (JSON Lines)

Each file contains:
1. **Meta line** (first line): Document metadata
2. **Chunk lines**: Individual chunks with embeddings

**Example:**
```jsonl
{"type": "meta", "document": "api/authentication.md", "hash": "abc123", "indexed_at": "2024-01-15T14:32:00Z"}
{"type": "chunk", "section": "JWT Token Validation", "offset": 1024, "length": 512, "embedding": [0.1, 0.2, ...]}
{"type": "chunk", "section": "Token Refresh", "offset": 2048, "length": 256, "embedding": [0.3, 0.4, ...]}
```

---

## Performance Considerations

### Indexing Performance

- **Small projects** (< 100 docs): ~1-5 seconds
- **Medium projects** (100-500 docs): ~10-30 seconds
- **Large projects** (500+ docs): ~1-5 minutes

Performance depends on:
- TEI server capacity
- Document complexity
- Chunk size settings
- Network latency

### Search Performance

- Typical response time: 50-200ms
- Primarily depends on TEI embedding generation
- `--top-k` has minimal impact on performance
- JSON output is slightly faster than formatted output

---

## Troubleshooting

### TEI Connection Issues

```bash
# Check if TEI is running
curl http://localhost:8080/health

# Check TEI logs
docker logs <container-id>  # If using Docker
```

### Index Not Updating

```bash
# Force reindex
doc-index index --force

# Check document hashes
doc-index status
```

### Poor Search Results

- Lower `--min-score` threshold
- Increase `--top-k` value
- Use more descriptive queries
- Check if documents are properly indexed