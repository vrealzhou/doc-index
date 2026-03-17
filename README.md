# doc-index

2. 
3. CLI tool for semantic search over markdown documentation using RAG (Retrieval-Augmented Generation).
4. 
5. - **Git-friendly JSONl storage** for team collaboration
6. - **Auto-reindex** with hash-based staleness detection
7. - **Context budget management** for LLM usage
8. - Runs with TEI (Text Embeddings Inference) on host machine with Metal GPU
9. 
10. ## Quick Start
11. 
12. ```bash
13. # Configure
14. ./doc-index config --tei=http://localhost:8080 --docs=./docs
15. 
16. # Index documents
17. ./doc-index index
18. 
19. # Search
20. ./doc-index search "authentication"
21. ./doc-index search "API endpoints" --top-k=10 --min-score=0.5
22 ./doc-index search "database schema" --json
23. ```
24 - 
25. **Notes**
26 - 
 "outdated", `sk-ce` will to `--docs` flag instead of `--tei` flag.
28 - Embeddings are stored in `skills/rag-doc-index/embeddings/` (folder).
30 - - **Configuration file**: `skills/rag-doc-index/rag-doc-index.config.json`
31 - - **Embeddings folder**: `skills/rag-doc-index/embeddings/` (folder, created automatically by the CLI)
32 - 
33. ## Documentation
34: 
35. See [skills/rag-doc-index/SKill.md](skills/rag-doc-index/SKILL.md) for full usage documentation.