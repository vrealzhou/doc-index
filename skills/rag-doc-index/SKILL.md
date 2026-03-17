---
name: rag-doc-index
description: >
  Use this skill when you need to search documentation in a codebase using semantic
  search. It indexes markdown documents and provides RAG (Retrieval-Augmented Generation)
  capabilities for finding relevant sections. Use it when the user asks about project
  documentation, design docs, architecture notes, or any markdown-based knowledge base.
---

# RAG Document Index Skill

Semantic search over markdown documentation using RAG (Retrieval-Augmented Generation).

## Overview

This skill enables intelligent search through markdown documentation by:
- **Indexing** markdown files and generating vector embeddings
- **Semantic search** using natural language queries
- **Context retrieval** for LLM-based workflows

## Sub-Skills

This skill is organized into two focused sub-skills:

### 📇 [Index Skill](./index/SKILL.md)
Set up and index your markdown documentation.
- Generate embeddings from markdown files
- Manage document updates and reindexing
- Check indexing status

**Use when:** You need to index or update the document index.

### 🔍 [Search Skill](./search/SKILL.md)
Search indexed documentation using semantic queries.
- Natural language search queries
- Filter and rank results by relevance
- Retrieve precise document sections

**Use when:** Documents are already indexed and you need to find relevant information.

## Prerequisites

### System Requirements

- **Go 1.26+** installed
- **TEI Server** (Text Embeddings Inference) for generating embeddings
- **Markdown documentation** files to index

### Hardware Recommendations

- **macOS**: Apple Silicon with Metal support
- **Linux**: NVIDIA GPU with CUDA support (recommended) or CPU
- **Memory**: At least 4GB RAM for embedding generation

## Installation

### Step 1: Install TEI Server

TEI (Text Embeddings Inference) is required for generating vector embeddings.

**macOS (Apple Silicon):**
```bash
text-embeddings-inference \
  --model-id BAAI/bge-small-en-v1.5 \
  --port 8080 \
  --device metal
```

**Linux (with NVIDIA GPU):**
```bash
docker run --gpus all -p 8080:80 \
  ghcr.io/huggingface/text-embeddings-inference:latest \
  --model-id BAAI/bge-small-en-v1.5
```

**Linux (CPU only):**
```bash
docker run -p 8080:80 \
  ghcr.io/huggingface/text-embeddings-inference:latest \
  --model-id BAAI/bge-small-en-v1.5
```

**Verify TEI is running:**
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

### Step 2: Install the CLI Tool

Install the `doc-index` CLI tool using Go:

```bash
go install github.com/vrealzhou/doc-index/cmd/doc-index@latest
```

Verify the installation:

```bash
doc-index --version
```

## Configuration

### Initial Setup

Configure the skill with your TEI endpoint and documentation path:

```bash
doc-index config --tei=http://localhost:8080 --docs=/path/to/project/docs
```

### Configuration Options

| Option | Description | Example |
|--------|-------------|---------|
| `--tei=<url>` | TEI server endpoint URL | `http://localhost:8080` |
| `--docs=<path>` | Root path to markdown documents | `/project/docs` |

### View Current Configuration

Check your current settings:

```bash
doc-index config --show
```

### Configuration Storage

Configuration is saved to: `skills/rag-doc-index/rag-doc-index.config.json`

```json
{
  "tei_endpoint": "http://localhost:8080",
  "docs_path": "/path/to/project/docs"
}
```

## Quick Start

```bash
# 1. Configure the system
doc-index config --tei=http://localhost:8080 --docs=/path/to/docs

# 2. Index your documents (see index/SKILL.md for details)
doc-index index

# 3. Search your documents (see search/SKILL.md for details)
doc-index search "authentication flow"
doc-index search "API endpoints" --top-k=10
```

## When to Use This Skill

Use this skill when you need to:
- Find relevant documentation sections in a codebase
- Understand project architecture or design decisions
- Search for specific concepts or features in documentation
- Retrieve context about APIs, configuration, or best practices

## Key Features

- ✅ **Git-friendly JSONl storage** for team collaboration
- ✅ **Auto-reindex** with hash-based staleness detection
- ✅ **Context budget management** for LLM usage
- ✅ **Metal GPU support** on macOS for fast embeddings
- ✅ **Semantic search** with natural language queries
- ✅ **Programmatic access** via JSON output

## Architecture

```
your-project/
├── docs/                    # Your markdown documentation
│   ├── README.md
│   ├── architecture/
│   └── api/
└── embeddings/              # Auto-generated embeddings
    ├── README.jsonl
    ├── architecture.jsonl
    └── api.jsonl
```

## Storage

- **Configuration**: `skills/rag-doc-index/rag-doc-index.config.json`
- **Embeddings**: Sibling `embeddings/` directory (auto-created)
- **Format**: JSONL with metadata and chunk embeddings

## Related Documentation

- **[Index Documentation](../../docs/rag-doc-index-installation.md)** - Detailed installation guide
- **[Configuration Guide](../../docs/rag-doc-index-configuration.md)** - Configuration options
- **[API Reference](../../docs/rag-doc-index-api-reference.md)** - Complete command reference
- **[Usage Examples](../../docs/rag-doc-index-usage-examples.md)** - Detailed examples

## Quick Reference

| Goal | Command | See |
|------|---------|-----|
| Configure | `doc-index config --tei=<url> --docs=<path>` | This file |
| Index documents | `doc-index index` | [Index Skill](./index/SKILL.md) |
| Search documents | `doc-index search "query"` | [Search Skill](./search/SKILL.md) |
| Check status | `doc-index status` | [Index Skill](./index/SKILL.md) |

## Getting Started Checklist

- [ ] Install TEI server and verify it's running
- [ ] Install doc-index CLI tool
- [ ] Configure with your docs path: `doc-index config --tei=http://localhost:8080 --docs=/path/to/docs`
- [ ] Index your documents: `doc-index index`
- [ ] Test search: `doc-index search "installation"`
- [ ] Explore results and adjust parameters as needed

## Next Steps

1. 📇 **[Index your documents](./index/SKILL.md)** - Learn how to index and manage your documentation
2. 🔍 **[Search your documents](./search/SKILL.md)** - Learn how to search and retrieve relevant information

## Notes

- TEI server must be running for both indexing and search operations
- Embeddings are stored in a sibling `embeddings/` directory to the docs folder
- The system uses hash-based detection to avoid reindexing unchanged documents
- Use natural language queries for best search results
- GPU acceleration significantly improves performance (Metal on macOS, CUDA on Linux)