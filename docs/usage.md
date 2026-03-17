# doc-index CLI

Command-line tool for semantic search over markdown documentation.

## Overview

doc-index indexes markdown documents and provides semantic search using vector embeddings. It supports multiple embedding providers including TEI (Text Embeddings Inference), oMLX (Apple Silicon optimized), and OpenAI.

## Installation

```bash
go build -o doc-index ./cmd/doc-index
```

## Quick Start

```bash
# Configure for your embedding provider
./doc-index config --provider=omlx --endpoint=http://localhost:8000 --api-key=your-key

# Index documents
./doc-index index

# Search
./doc-index search "authentication"
```

## Commands

### config

Configure embedding provider and settings.

```bash
./doc-index config [flags]
```

**Flags:**
- `--provider`: Embedding provider (`tei`, `omlx`, `openai`)
- `--endpoint`: Embedding service URL
- `--api-key`: API key for authentication
- `--model`: Embedding model name
- `--docs`: Documents root path
- `--show`: Display current configuration
- `--skill-folder`: Override skill folder location

**Examples:**
```bash
# Configure for TEI
./doc-index config --provider=tei --endpoint=http://localhost:8080

# Configure for oMLX with API key
./doc-index config --provider=omlx --endpoint=http://localhost:8000 --api-key=sk-xxx --model=bge-m3

# Configure for OpenAI
./doc-index config --provider=openai --endpoint=https://api.openai.com --api-key=sk-xxx --model=text-embedding-3-small

# Show current config
./doc-index config --show

# Update just the API key
./doc-index config --api-key=new-key
```

### index

Index all markdown documents in the configured docs path.

```bash
./doc-index index [flags]
```

**Flags:**
- `-f, --force`: Force reindex all documents
- `--skill-folder`: Override skill folder location

**Examples:**
```bash
# Index only changed documents
./doc-index index

# Force reindex all
./doc-index index --force
```

### search

Search for relevant documentation sections.

```bash
./doc-index search <query> [flags]
```

**Flags:**
- `-k, --top-k`: Number of results (default: 5)
- `-m, --min-score`: Minimum similarity score (default: 0.3)
- `-j, --json`: Output as JSON

**Examples:**
```bash
./doc-index search "authentication"
./doc-index search "API rate limiting" --top-k=10 --min-score=0.5
./doc-index search "database schema" --json
```

### status

Show index status and provider connectivity.

```bash
./doc-index status [flags]
```

**Flags:**
- `--skill-folder`: Override skill folder location

## Configuration

Configuration is stored in `rag-doc-index.config.json` within the skill folder:

```json
{
  "provider": "omlx",
  "endpoint": "http://localhost:8000",
  "api_key": "your-api-key",
  "model": "bge-m3",
  "docs_path": "docs"
}
```

### Skill Folder Locations

The tool searches for the skill folder in this order:
1. `.opencode/skills/rag-doc-index`
2. `.claude/skills/rag-doc-index`
3. `skills/rag-doc-index`
4. Current working directory (creates if not found)

Override with `--skill-folder` flag or `RAG_SKILL_FOLDER` environment variable.

## Supported Providers

| Provider | Endpoint Format | Authentication |
|----------|----------------|----------------|
| `tei` | `http://host:port` | Optional Bearer token |
| `omlx` | `http://host:port` | Bearer token (OpenAI-compatible) |
| `openai` | `https://api.openai.com` | Bearer token required |

## Storage

- **Config**: `skills/rag-doc-index/rag-doc-index.config.json`
- **Embeddings**: `skills/rag-doc-index/embeddings/*.jsonl`
- **Format**: JSONL with one meta line + chunk lines per document

## Environment Variables

All config options can be set via environment variables:

- `EMBEDDING_PROVIDER`: Provider type (tei, omlx, openai)
- `EMBEDDING_ENDPOINT`: Service endpoint URL
- `API_KEY`: API key for authentication
- `EMBEDDING_MODEL`: Model name
- `DOCS_PATH`: Documents root path
- `EMBEDDINGS_PATH`: Embeddings storage path
- `RAG_SKILL_FOLDER`: Override skill folder location

## Git Integration

Add to `.gitattributes` for team collaboration:

```
embeddings/*.jsonl merge=union
```

## Requirements

- Go 1.26+
- Embedding service (TEI, oMLX, or OpenAI account)
- macOS 15+ for oMLX (Apple Silicon)
