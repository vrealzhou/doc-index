---
name: rag-doc-index
description: >
  Use this skill when you need to search documentation in a codebase using semantic
  search. It indexes markdown documents and provides RAG (Retrieval-Augmented Generation)
  capabilities for finding relevant sections. Use it when the user asks about project
  documentation, design docs, architecture notes, or any markdown-based knowledge base.
---

# RAG Document Index Skill

Semantic search over markdown documentation for coding agents.

## Prerequisites

This skill requires a running TEI (Text Embeddings Inference) server.

### Start TEI Server

**macOS (Apple Silicon):**
```bash
text-embeddings-inference \
  --model-id BAAI/bge-small-en-v1.5 \
  --port 8080 \
  --device metal
```

**Linux (with GPU):**
```bash
docker run --gpus all -p 8080:80 \
  ghcr.io/huggingface/text-embeddings-inference:latest \
  --model-id BAAI/bge-small-en-v1.5
```

Verify TEI is running:
```bash
curl http://localhost:8080/health
```

## Installation

### Step 1 — Install the CLI tool
need go v1.26+
```bash
go install github.com/vrealzhou/doc-index/cmd/doc-index@latest
```

## Step 2 — Configure the skill

Set the TEI endpoint and documents path:

```bash
doc-index config --tei=http://localhost:8080 --docs=/path/to/project/docs
```

Show current configuration:
```bash
doc-index config --show
```

## Step 3 — Index documents

Index all markdown files in the documents path:

```bash
doc-index index
```

Force reindex all documents:
```bash
doc-index index --force
```

## Step 4 — Search for relevant documentation

Search using natural language queries:

```bash
doc-index search "authentication flow"
doc-index search "API rate limiting" --top-k=10 --min-score=0.5
```

Get JSON output for programmatic use:
```bash
doc-index search "database schema" --json
```

## Step 5 — Use results in your work

The search results include:
- **Document**: Source markdown file
- **Section**: Heading/title of the relevant section
- **Position**: Offset and length for reading specific content
- **Score**: Similarity score (0-1, higher is better)

Use these results to understand the codebase context before making changes.

## Quick Reference

| Goal | Command |
|------|---------|
| Configure | `doc-index config --tei=<url> --docs=<path>` |
| Show config | `doc-index config --show` |
| Index documents | `doc-index index` |
| Force reindex | `doc-index index --force` |
| Check status | `doc-index status` |
| Basic search | `doc-index search "query"` |
| More results | `doc-index search "query" --top-k=10` |
| Filter by score | `doc-index search "query" --min-score=0.5` |
| JSON output | `doc-index search "query" --json` |

## Options

### search command

| Option | Default | Description |
|--------|---------|-------------|
| `--top-k=N` | 5 | Number of results to return |
| `--min-score=N` | 0.3 | Minimum similarity score (0-1) |
| `--json` | false | Output as JSON |

### config command

| Option | Description |
|--------|-------------|
| `--tei=<url>` | Set TEI endpoint URL |
| `--docs=<path>` | Set documents root path |
| `--show` | Show current configuration |

## Configuration File

Configuration is saved to `rag-doc-index.config.json` alongside the binary:

```json
{
  "tei_endpoint": "http://localhost:8080",
  "docs_path": "/path/to/project/docs"
}
```

## Storage

- **Config**: `skills/rag-doc-index/rag-doc-index.config.json`
- **Embeddings**: `<docs_path>/../embeddings/*.jsonl`
- **Format**: JSONL with meta line + chunk lines

## Notes

- Embeddings are stored in a sibling `embeddings/` directory to the docs folder
- Use `--force` to reindex when documents have been modified
- Lower `--min-score` if getting no results (default: 0.3)
- TEI must be running for indexing and search operations
