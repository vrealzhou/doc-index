---
name: rag-doc-index-index
description: >
  Use this skill when you need to index markdown documentation files for semantic
  search. This skill handles indexing documents, managing embeddings, and checking
  indexing status. Use it when documents need to be indexed for the first time or
  reindexed after updates. Requires prior setup (see main rag-doc-index skill).
---

# RAG Document Indexing

Index markdown documentation files for semantic search using vector embeddings.

## Overview

This skill handles indexing operations for the RAG Document Index system:
- Indexing markdown documents
- Managing embeddings storage
- Checking indexing status
- Updating the index when documents change

## Prerequisites

Before using this skill, ensure you have:
- TEI server running on localhost:8080
- doc-index CLI tool installed
- Configuration set up with `doc-index config`

**Setup Required**: If not yet installed, see the main `rag-doc-index` skill for installation and configuration instructions.

## Basic Usage

### Index Documents

Index all markdown files in the configured path:

```bash
doc-index index
```

This will:
- Scan all `.md` files in your docs directory
- Generate embeddings using TEI
- Store embeddings in JSONL format
- Skip documents that haven't changed (hash-based detection)

### Force Reindex

Reindex all documents, even if they haven't changed:

```bash
doc-index index --force
```

Use when:
- Documents have been modified
- TEI model has been updated
- Embeddings seem corrupted or outdated
- After upgrading the CLI tool

### Check Indexing Status

View current indexing status and statistics:

```bash
doc-index status
```

Output includes:
- Number of documents indexed
- Total chunks created
- Last indexing timestamp
- Embeddings location
- TEI server status

## Indexing Options

| Option | Type | Description |
|--------|------|-------------|
| `--force` | flag | Force reindex all documents |

## Understanding Indexing

### What Gets Indexed

- All `.md` files in the configured docs path
- Documents are chunked by headers and sections
- Each chunk gets a vector embedding
- Metadata is stored for efficient retrieval

### Hash-Based Detection

The system uses hash-based staleness detection:
- Skips unchanged documents (faster reindexing)
- Detects modified documents automatically
- Use `--force` to override and reindex everything

### Storage Location

Embeddings are stored in:
- **Path**: `<docs_path>/../embeddings/*.jsonl`
- **Format**: JSONL with metadata and chunk embeddings
- **Auto-created**: Directory is created automatically

**Example Structure:**
```
your-project/
├── docs/                    # Documentation root
│   ├── README.md
│   ├── architecture/
│   └── api/
└── embeddings/              # Auto-generated
    ├── README.jsonl
    ├── architecture.jsonl
    └── api.jsonl
```

### JSONL Format

Each file contains:
1. **Meta line**: Document metadata (path, hash, timestamp)
2. **Chunk lines**: Individual chunks with embeddings

**Example:**
```jsonl
{"type": "meta", "document": "api/authentication.md", "hash": "abc123", "indexed_at": "2024-01-15T14:32:00Z"}
{"type": "chunk", "section": "JWT Token Validation", "offset": 1024, "length": 512, "embedding": [0.1, 0.2, ...]}
{"type": "chunk", "section": "Token Refresh", "offset": 2048, "length": 256, "embedding": [0.3, 0.4, ...]}
```

## Common Workflows

### Initial Indexing

```bash
# After configuration (see main skill)
doc-index index

# Verify indexing
doc-index status
```

### After Documentation Updates

```bash
# Check current status
doc-index status

# Index new/modified documents
doc-index index

# Or force reindex everything
doc-index index --force
```

### Troubleshooting Indexing

```bash
# Check if documents are indexed
doc-index status

# Verify configuration
doc-index config --show

# Force complete reindex
doc-index index --force
```

## Troubleshooting

### No Documents Found

**Problem**: Indexing finds no documents

**Solutions**:
1. Verify docs path: `doc-index config --show`
2. Check that `.md` files exist in the directory
3. Ensure read permissions for the docs directory

### Indexing Fails

**Problem**: Indexing fails with errors

**Solutions**:
1. Verify TEI is running: `curl http://localhost:8080/health`
2. Check TEI logs for errors
3. Ensure sufficient memory (at least 4GB)
4. Try force reindex: `doc-index index --force`

### Slow Indexing

**Problem**: Indexing is very slow

**Solutions**:
1. Use GPU acceleration if available
2. Ensure TEI is running on GPU
3. Reduce document size or complexity
4. Check available memory

### Outdated Index

**Problem**: Index doesn't reflect recent changes

**Solutions**:
1. Check status: `doc-index status`
2. Run incremental index: `doc-index index`
3. Or force reindex: `doc-index index --force`

## Performance

### Expected Performance

- **Small projects** (< 100 docs): ~1-5 seconds
- **Medium projects** (100-500 docs): ~10-30 seconds
- **Large projects** (500+ docs): ~1-5 minutes

Performance depends on:
- TEI server capacity
- Document complexity
- Available memory
- GPU acceleration

## Quick Reference

| Goal | Command |
|------|---------|
| Index documents | `doc-index index` |
| Force reindex | `doc-index index --force` |
| Check status | `doc-index status` |

## Related Skills

- **rag-doc-index** (parent): Installation and setup instructions
- **rag-doc-index-search**: Search indexed documents

## Notes

- Documents must be indexed before searching
- Use hash-based detection for faster reindexing
- TEI must be running for indexing operations
- Embeddings are stored in a sibling `embeddings/` directory