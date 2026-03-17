# Installation Guide

This guide covers the installation and setup of the RAG Document Index skill.

## Prerequisites

### System Requirements

- **Go** v1.26 or later (for building the CLI tool)
- **TEI Server** (Text Embeddings Inference) for generating embeddings
- **Markdown documentation** files to index

### Hardware Recommendations

- **macOS**: Apple Silicon with Metal support
- **Linux**: NVIDIA GPU with CUDA support (recommended) or CPU
- **Memory**: At least 4GB RAM for embedding generation

## Step 1: Install TEI Server

TEI (Text Embeddings Inference) is required for generating vector embeddings.

### macOS (Apple Silicon)

Install and run TEI with Metal GPU acceleration:

```bash
# Install TEI (if not already installed)
# Follow instructions at: https://github.com/huggingface/text-embeddings-inference

# Start TEI server
text-embeddings-inference \
  --model-id BAAI/bge-small-en-v1.5 \
  --port 8080 \
  --device metal
```

### Linux (with NVIDIA GPU)

Run TEI using Docker with GPU support:

```bash
docker run --gpus all -p 8080:80 \
  ghcr.io/huggingface/text-embeddings-inference:latest \
  --model-id BAAI/bge-small-en-v1.5
```

### Linux (CPU only)

```bash
docker run -p 8080:80 \
  ghcr.io/huggingface/text-embeddings-inference:latest \
  --model-id BAAI/bge-small-en-v1.5
```

### Verify TEI Installation

Check that TEI is running properly:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"ok"}
```

## Step 2: Install the CLI Tool

Install the `doc-index` CLI tool using Go:

```bash
go install github.com/vrealzhou/doc-index/cmd/doc-index@latest
```

Verify the installation:

```bash
doc-index --version
```

## Step 3: Initial Configuration

Configure the CLI tool with your TEI endpoint and documentation path:

```bash
doc-index config \
  --tei=http://localhost:8080 \
  --docs=/path/to/your/project/docs
```

### Configuration Options

| Option | Description | Example |
|--------|-------------|---------|
| `--tei` | TEI server endpoint | `http://localhost:8080` |
| `--docs` | Path to markdown documentation | `/project/docs` |

### View Current Configuration

Check your current settings:

```bash
doc-index config --show
```

## Step 4: Index Your Documents

Create the initial index of your markdown documentation:

```bash
doc-index index
```

This will:
- Scan all `.md` files in your docs directory
- Generate embeddings using TEI
- Store embeddings in JSONL format
- Create metadata for efficient retrieval

### Verify Indexing

Check the index status:

```bash
doc-index status
```

## Step 5: Test Search Functionality

Verify everything is working with a test search:

```bash
doc-index search "installation" --top-k=3
```

## Troubleshooting

### TEI Server Issues

**Problem**: Cannot connect to TEI server

**Solutions**:
- Verify TEI is running: `curl http://localhost:8080/health`
- Check firewall settings
- Ensure correct port (default: 8080)

### Indexing Issues

**Problem**: No documents found

**Solutions**:
- Verify docs path is correct: `doc-index config --show`
- Check that markdown files exist in the specified directory
- Ensure read permissions for the docs directory

**Problem**: Embedding generation fails

**Solutions**:
- Verify TEI server is running
- Check TEI logs for errors
- Ensure sufficient memory is available

### Performance Issues

**Problem**: Slow indexing or search

**Solutions**:
- Use GPU acceleration if available
- Reduce the number of concurrent requests to TEI
- Consider using a smaller embedding model

## Next Steps

After successful installation:

1. Review the [Configuration Guide](./configuration.md) for advanced options
2. See [Basic Usage](../examples/basic-usage.md) for common use cases
3. Check the [API Reference](./api-reference.md) for all available commands

## Upgrading

To upgrade to the latest version:

```bash
go install github.com/vrealzhou/doc-index/cmd/doc-index@latest
```

After upgrading, you may need to reindex your documents:

```bash
doc-index index --force
```

## Uninstallation

To remove the CLI tool:

```bash
# Remove binary
rm $(which doc-index)

# Optional: Remove configuration and embeddings
rm -rf skills/rag-doc-index/
```

Stop the TEI server if no longer needed.