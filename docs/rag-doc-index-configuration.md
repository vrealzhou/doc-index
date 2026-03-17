# Configuration

Configuration for the RAG Document Index skill is managed through CLI commands and stored in a JSON file.

## Quick Configuration

Set up the skill with the TEI endpoint and documents path:

```bash
doc-index config --tei=http://localhost:8080 --docs=/path/to/project/docs
```

## View Current Configuration

Display the current configuration settings:

```bash
doc-index config --show
```

## Configuration Options

| Option | Description | Example |
|--------|-------------|---------|
| `--tei=<url>` | Text Embeddings Inference server endpoint | `--tei=http://localhost:8080` |
| `--docs=<path>` | Root path to markdown documentation | `--docs=./docs` |
| `--show` | Display current configuration | `--show` |

## Configuration File

The configuration is automatically saved to:

```
skills/rag-doc-index/rag-doc-index.config.json
```

### File Structure

```json
{
  "tei_endpoint": "http://localhost:8080",
  "docs_path": "/path/to/project/docs"
}
```

### Fields

| Field | Type | Description |
|-------|------|-------------|
| `tei_endpoint` | string | URL of the TEI server (Text Embeddings Inference) |
| `docs_path` | string | Absolute or relative path to the documentation root directory |

## Storage Locations

The skill uses several storage locations:

| Type | Path | Description |
|------|------|-------------|
| **Configuration** | `skills/rag-doc-index/rag-doc-index.config.json` | Skill configuration file |
| **Embeddings** | `<docs_path>/../embeddings/*.jsonl` | Embeddings storage (sibling to docs folder) |
| **Index** | `<docs_path>/../embeddings/` | JSONL files with metadata and chunks |

## Embeddings Storage Format

Embeddings are stored in JSONL format with:
- **Meta line**: Document metadata (path, hash, timestamp)
- **Chunk lines**: Individual text chunks with embeddings

### Example Structure

```
your-project/
├── docs/                    # Documentation root (--docs path)
│   ├── README.md
│   ├── architecture/
│   └── api/
└── embeddings/              # Auto-generated sibling folder
    ├── README.jsonl        # Embeddings for README.md
    ├── architecture.jsonl  # Embeddings for architecture/
    └── api.jsonl           # Embeddings for api/
```

## Environment Requirements

### TEI Server

The TEI (Text Embeddings Inference) server must be running and accessible:

- **Endpoint**: Configured via `--tei` flag
- **Health Check**: `curl http://localhost:8080/health`
- **Required**: Must be running for both indexing and search operations

### Documents Path

The documents path should:
- Contain markdown files (`.md` extension)
- Be readable by the CLI tool
- Be a directory (not a file)

## Common Configuration Patterns

### Local Development

```bash
doc-index config \
  --tei=http://localhost:8080 \
  --docs=./docs
```

### Project with Nested Documentation

```bash
doc-index config \
  --tei=http://localhost:8080 \
  --docs=/projects/my-app/documentation
```

### Multiple Projects

Create separate configuration files by running the config command in each project directory:

```bash
# Project A
cd /projects/project-a
doc-index config --tei=http://localhost:8080 --docs=./docs

# Project B
cd /projects/project-b
doc-index config --tei=http://localhost:8080 --docs=./documentation
```

## Troubleshooting

### Configuration Not Found

If you see configuration errors:
1. Run `doc-index config --show` to check current settings
2. Re-run the config command with proper paths
3. Ensure the configuration file is writable

### TEI Connection Issues

If the TEI server is unreachable:
1. Verify TEI is running: `curl http://localhost:8080/health`
2. Check the endpoint URL in configuration
3. Ensure no firewall is blocking the connection

### Documents Path Issues

If documents aren't being found:
1. Verify the path exists: `ls -la /path/to/docs`
2. Check the path contains `.md` files
3. Ensure the path is correctly configured with `--docs`
