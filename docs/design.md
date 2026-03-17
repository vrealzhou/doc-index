# doc-index Design Document

## Overview

doc-index is a CLI tool for semantic search over markdown documentation. It indexes documents using vector embeddings and provides progressive disclosure of content to coding agents without overwhelming context windows.

## Architecture

### Components

```
cmd/doc-index/
├── main.go              # Entry point
└── cmd/
    ├── root.go          # Root command (cobra)
    ├── config.go        # Configuration management
    ├── index.go         # Document indexing
    ├── search.go        # Semantic search
    └── status.go        # Status reporting

internal/
├── config/              # Configuration loading
├── chunk/               # Document chunking
├── embed/               # Embedding client (multi-provider)
├── indexer/             # Index management
└── search/              # Search engine

skills/rag-doc-index/
└── SKILL.md             # Skill documentation
```

## Data Flow

### Indexing Flow

```
Markdown Files → Chunker → Embedder → JSONL Storage
     ↓              ↓          ↓            ↓
   Read         Split by    Call API    Save to
   Content      Headers     (TEI/       skill/
                            oMLX/       embeddings/
                            OpenAI)
```

1. **Scan**: Walk docs directory recursively
2. **Chunk**: Split by markdown headers with overlap
3. **Embed**: Call embedding provider API
4. **Store**: Save as JSONL in skill folder

### Search Flow

```
Query → Embedder → Vector Search → Rank → Return Results
  ↓        ↓            ↓           ↓          ↓
Text    Get         Cosine      Position   Doc ID,
        Vector      Similarity   Aware      Score,
                    (brute                      Preview
                    force)
```

1. **Embed Query**: Convert query to vector
2. **Search**: Brute-force cosine similarity
3. **Rank**: Position-aware ordering
4. **Return**: Document references with scores

## Storage Format

### JSONL Structure

Each document has a corresponding `.jsonl` file:

```jsonl
{"meta":true,"v":1,"doc":"api-design.md","hash":"a1b2c3d4","chunks":3,"mtime":1700000000}
{"idx":0,"title":"Overview","offset":0,"length":500,"vec":[0.123,-0.456,...]}
{"idx":1,"title":"Authentication","offset":500,"length":450,"vec":[0.789,0.234,...]}
```

**Meta line fields:**
- `meta`: Always true for meta line
- `v`: Format version
- `doc`: Relative document path (supports subfolders)
- `hash`: SHA256 content hash (first 16 chars)
- `chunks`: Number of chunks
- `mtime`: Index timestamp

**Chunk line fields:**
- `idx`: Chunk index
- `title`: Section title
- `offset`: Character offset in document
- `length`: Content length
- `vec`: Embedding vector (384 dims for BGE-small)

### Directory Structure

```
repo-root/
├── docs/                          # Source markdown files
│   ├── api-design.md
│   └── guides/
│       └── setup.md
├── skills/rag-doc-index/          # Skill folder
│   ├── rag-doc-index.config.json  # Configuration
│   └── embeddings/                # Generated embeddings
│       ├── api-design.md.jsonl
│       └── guides/
│           └── setup.md.jsonl
```

## Configuration

### File Config (rag-doc-index.config.json)

```json
{
  "provider": "omlx",
  "endpoint": "http://localhost:8000",
  "api_key": "sk-xxx",
  "model": "bge-m3",
  "docs_path": "docs"
}
```

**Fields:**
- `provider`: Embedding provider (`tei`, `omlx`, `openai`)
- `endpoint`: Service URL
- `api_key`: Authentication token
- `model`: Model name
- `docs_path`: Documents root (relative to repo)

### Environment Variables

- `EMBEDDING_PROVIDER`: Provider type
- `EMBEDDING_ENDPOINT`: Service URL
- `API_KEY`: Authentication
- `EMBEDDING_MODEL`: Model name
- `DOCS_PATH`: Documents path
- `EMBEDDINGS_PATH`: Embeddings path
- `RAG_SKILL_FOLDER`: Override skill location

### Precedence

1. Command-line flags
2. Environment variables
3. Config file
4. Defaults

## Embedding Providers

### TEI (Text Embeddings Inference)

- **Endpoint**: `http://host:port`
- **API**: `/embed` (HuggingFace format)
- **Auth**: Optional Bearer token
- **Best for**: Self-hosted, GPU servers

### oMLX

- **Endpoint**: `http://host:port`
- **API**: `/v1/embeddings` (OpenAI-compatible)
- **Auth**: Bearer token
- **Best for**: Apple Silicon (M1/M2/M3/M4)
- **Features**: Tiered KV cache, continuous batching

### OpenAI

- **Endpoint**: `https://api.openai.com`
- **API**: `/v1/embeddings`
- **Auth**: Bearer token required
- **Best for**: Cloud API, no local setup

## Chunking Strategy

### Algorithm

1. **Split by Headers**: Use `## ` and `# ` as section boundaries
2. **Extract Title**: Header text becomes chunk title
3. **Overlap**: 100 characters overlap between chunks
4. **Max Size**: 1200 characters per chunk

### Example

```markdown
# Main Title

Intro text here.

## Section One

Content for section one.
More content here that might be long.

## Section Two

Content for section two.
```

**Chunks:**
1. Title: "Overview", Content: "# Main Title\n\nIntro text here."
2. Title: "Section One", Content: "## Section One\n\nContent..."
3. Title: "Section Two", Content: "## Section Two\n\nContent..."

## Search Algorithm

### Cosine Similarity

```go
similarity = dot(query, doc) / (norm(query) * norm(doc))
```

### Position-Aware Ordering

When multiple chunks from the same document match:
1. Sort by index
2. Return consecutive chunks together
3. Improves context coherence

### Filtering

- **Min Score**: Default 0.3 (configurable)
- **Top K**: Default 5 results (configurable)

## Context Budget Management

### Constraints

- **Total Context**: ~200KB for Claude Code
- **Search Budget**: ~100KB (50% of total)
- **Per Result**: ~150 chars preview

### Calculation

```
Tokens ≈ Characters / 4
100KB ≈ 25K tokens
```

### Progressive Disclosure

1. **Search**: Returns previews only
2. **Read**: Fetch full sections on demand
3. **Outline**: Get document structure

## Git Integration

### Merge Strategy

Add to `.gitattributes`:

```
embeddings/*.jsonl merge=union
```

This allows concurrent indexing by multiple developers.

### Staleness Detection

- Compare content hash (SHA256)
- Reindex when hash changes
- Remove orphaned embeddings

## Security

### API Keys

- Stored in config file (user-readable only)
- Masked in status output (show first/last 4 chars)
- Passed as Bearer token in Authorization header

### File Permissions

- Config: `0644` (user read/write)
- Embeddings: `0755` directories, `0644` files

## Performance

### Indexing

- **Batch Size**: 32 texts per API call
- **Timeout**: 30 seconds per request
- **Concurrent**: Sequential (respect API limits)

### Search

- **Algorithm**: Brute-force cosine similarity
- **Complexity**: O(N) where N = total chunks
- **Typical**: <10ms for 2000 chunks
- **Scale**: Suitable for <100 documents

### Memory

- **Embeddings**: ~3MB per 1000 chunks (384 dims)
- **Index**: Loaded entirely in memory
- **Search**: No disk I/O during query

## Future Enhancements

### Potential Improvements

1. **HNSW Index**: For larger document sets (>10K chunks)
2. **Incremental Updates**: Only changed sections
3. **Multi-language**: Support for non-English docs
4. **Caching**: LRU cache for frequent queries
5. **Compression**: Quantized embeddings (int8)

### Out of Scope

- Real-time collaboration
- Distributed indexing
- Web interface
- Cloud storage backends
