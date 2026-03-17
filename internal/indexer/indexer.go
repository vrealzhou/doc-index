package indexer

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vrealzhou/doc-index/internal/chunk"
	"github.com/vrealzhou/doc-index/internal/config"
	"github.com/vrealzhou/doc-index/internal/embed"
)

type DocStatus int

const (
	StatusCurrent DocStatus = iota
	StatusStale
	StatusMissing
	StatusOrphan
)

type StatusReport struct {
	Total   int                  `json:"total"`
	Current int                  `json:"current"`
	Stale   int                  `json:"stale"`
	Missing int                  `json:"missing"`
	Orphan  int                  `json:"orphan"`
	Details map[string]DocStatus `json:"details,omitempty"`
}

type Indexer struct {
	cfg      config.Config
	embedder *embed.Client
	mu       sync.RWMutex
}

func New(cfg config.Config, embedder *embed.Client) *Indexer {
	return &Indexer{
		cfg:      cfg,
		embedder: embedder,
	}
}

func (i *Indexer) Scan() (*StatusReport, error) {
	docsDir := i.cfg.DocsPath
	embedDir := i.cfg.EmbeddingsPath
	docsBase := filepath.Base(docsDir)

	report := &StatusReport{
		Details: make(map[string]DocStatus),
	}

	embedHashes := make(map[string]string)
	collectEmbedHashes(embedDir, embedDir, embedHashes)

	err := filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		relPath, err := filepath.Rel(docsDir, path)
		if err != nil {
			return err
		}

		docPath := filepath.Join(docsBase, relPath)

		report.Total++
		content, err := os.ReadFile(path)
		if err != nil {
			report.Details[docPath] = StatusMissing
			report.Missing++
			return nil
		}

		contentHash := computeHash(content)
		storedHash, exists := embedHashes[docPath]

		if !exists {
			report.Details[docPath] = StatusMissing
			report.Missing++
		} else if storedHash == contentHash {
			report.Details[docPath] = StatusCurrent
			report.Current++
		} else {
			report.Details[docPath] = StatusStale
			report.Stale++
		}
		delete(embedHashes, docPath)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk docs dir: %w", err)
	}

	for docName := range embedHashes {
		report.Details[docName] = StatusOrphan
		report.Orphan++
	}

	return report, nil
}

func collectEmbedHashes(baseDir, currentDir string, hashes map[string]string) {
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return
	}

	for _, e := range entries {
		fullPath := filepath.Join(currentDir, e.Name())
		if e.IsDir() {
			collectEmbedHashes(baseDir, fullPath, hashes)
			continue
		}
		if filepath.Ext(e.Name()) != ".jsonl" {
			continue
		}

		relPath, err := filepath.Rel(baseDir, fullPath)
		if err != nil {
			continue
		}

		relPath = strings.TrimSuffix(relPath, ".jsonl")
		hash, err := readStoredHash(fullPath)
		if err == nil {
			hashes[relPath] = hash
		}
	}
}

func (i *Indexer) Reindex(ctx context.Context, docNames []string) error {
	docsBase := filepath.Base(i.cfg.DocsPath)

	for _, name := range docNames {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		relPath := name
		if strings.HasPrefix(name, docsBase+"/") || strings.HasPrefix(name, docsBase+"\\") {
			relPath = name[len(docsBase)+1:]
		}

		docPath := filepath.Join(i.cfg.DocsPath, relPath)
		embedPath := filepath.Join(i.cfg.EmbeddingsPath, name+".jsonl")

		content, err := os.ReadFile(docPath)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}

		result := chunk.Split(string(content), name, i.cfg.MaxChunkSize, i.cfg.OverlapSize)

		texts := make([]string, len(result.Chunks))
		for j, c := range result.Chunks {
			start := c.Offset
			end := c.Offset + c.Length
			if end > len(content) {
				end = len(content)
			}
			texts[j] = string(content[start:end])
		}

		vectors, err := i.embedder.EmbedBatch(ctx, texts)
		if err != nil {
			return fmt.Errorf("embed %s: %w", name, err)
		}

		for j := range result.Chunks {
			result.Chunks[j].Vec = vectors[j]
		}

		if err := i.writeJSONL(embedPath, result); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}

	return nil
}

func (i *Indexer) AutoReindex(ctx context.Context) (*StatusReport, error) {
	report, err := i.Scan()
	if err != nil {
		return nil, err
	}

	var toReindex []string
	for name, status := range report.Details {
		if status == StatusStale || status == StatusMissing {
			toReindex = append(toReindex, name)
		}
	}

	for name, status := range report.Details {
		if status == StatusOrphan {
			embedPath := filepath.Join(i.cfg.EmbeddingsPath, name+".jsonl")
			os.Remove(embedPath)
		}
	}

	if len(toReindex) == 0 {
		return report, nil
	}

	if err := i.Reindex(ctx, toReindex); err != nil {
		return nil, err
	}

	return i.Scan()
}

func (i *Indexer) writeJSONL(path string, result chunk.ChunkWithMeta) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	tmpPath := path + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer os.Remove(tmpPath)

	writer := bufio.NewWriter(f)

	result.Meta.Mtime = time.Now().Unix()
	metaLine, _ := json.Marshal(result.Meta)
	writer.Write(metaLine)
	writer.WriteByte('\n')

	for _, c := range result.Chunks {
		line, _ := json.Marshal(c)
		writer.Write(line)
		writer.WriteByte('\n')
	}

	if err := writer.Flush(); err != nil {
		f.Close()
		return err
	}
	f.Close()

	return os.Rename(tmpPath, path)
}

func computeHash(content []byte) string {
	h := sha256.Sum256(content)
	return hex.EncodeToString(h[:])[:16]
}

func readStoredHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return "", fmt.Errorf("empty file")
	}

	var meta chunk.Meta
	if err := json.Unmarshal(scanner.Bytes(), &meta); err != nil {
		return "", err
	}

	return meta.Hash, nil
}
