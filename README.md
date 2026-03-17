# doc-index

CLI tool for semantic search over markdown documentation using RAG (Retrieval-Augmented Generation).

## Features

- **Git-friendly JSONl storage** for team collaboration
- **Auto-reindex** with hash-based staleness detection
- **Context budget management** for LLM usage
- Runs with TEI (Text Embeddings Inference) on host machine with Metal GPU

## Quick Start

```bash
# Configure
./doc-index config --tei=http://localhost:8080 --docs=./docs

# Index documents
./doc-index index

# Search
./doc-index search "authentication"
./doc-index search "API endpoints" --top-k=10 --min-score=0.5
./doc-index search "database schema" --json
```

## Notes

- If embeddings are marked as "outdated", `sk-ce` will refer to the `--docs` flag instead of `--tei` flag
- Embeddings are stored in `skills/rag-doc-index/embeddings/` folder
- **Configuration file**: `skills/rag-doc-index/rag-doc-index.config.json`
- **Embeddings folder**: `skills/rag-doc-index/embeddings/` (created automatically by the CLI)

## Documentation

See [skills/rag-doc-index/SKILL.md](skills/rag-doc-index/SKILL.md) for full usage documentation.