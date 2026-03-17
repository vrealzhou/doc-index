package search

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/vrealzhou/doc-index/internal/config"
	"github.com/vrealzhou/doc-index/internal/embed"
)

type Result struct {
	DocID   string  `json:"doc_id"`
	ChunkID int     `json:"chunk_id"`
	Title   string  `json:"title"`
	Score   float32 `json:"score"`
	Preview string  `json:"preview"`
	Offset  int     `json:"offset"`
	Length  int     `json:"length"`
}

type Document struct {
	ID     string
	Hash   string
	Chunks []ChunkData
}

type ChunkData struct {
	Idx    int
	Title  string
	Offset int
	Length int
	Vec    []float32
}

type Engine struct {
	cfg      config.Config
	embedder *embed.Client
	mu       sync.RWMutex
	docs     map[string]*Document
}

func NewEngine(cfg config.Config, embedder *embed.Client) *Engine {
	return &Engine{
		cfg:      cfg,
		embedder: embedder,
		docs:     make(map[string]*Document),
	}
}

func (e *Engine) LoadFromDisk(dir string) error {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".jsonl" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for _, path := range files {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			doc, err := e.loadSingleFile(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to load %s: %v\n", p, err)
				return
			}

			mu.Lock()
			e.docs[doc.ID] = doc
			mu.Unlock()
		}(path)
	}

	wg.Wait()
	return nil
}

func (e *Engine) loadSingleFile(path string) (*Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	var doc Document
	var chunks []ChunkData

	for scanner.Scan() {
		line := scanner.Bytes()

		var meta struct {
			IsMeta bool   `json:"meta"`
			DocID  string `json:"doc"`
			Hash   string `json:"hash"`
		}
		if err := json.Unmarshal(line, &meta); err == nil && meta.IsMeta {
			doc.ID = meta.DocID
			doc.Hash = meta.Hash
			continue
		}

		var c ChunkData
		var raw struct {
			Idx    int       `json:"idx"`
			Title  string    `json:"title"`
			Offset int       `json:"offset"`
			Length int       `json:"length"`
			Vec    []float32 `json:"vec"`
		}
		if err := json.Unmarshal(line, &raw); err == nil && raw.Vec != nil {
			c.Idx = raw.Idx
			c.Title = raw.Title
			c.Offset = raw.Offset
			c.Length = raw.Length
			c.Vec = raw.Vec
			chunks = append(chunks, c)
		}
	}

	doc.Chunks = chunks
	return &doc, scanner.Err()
}

func (e *Engine) Search(ctx context.Context, query string, topK int) ([]Result, error) {
	if topK <= 0 {
		topK = e.cfg.DefaultTopK
	}

	queryVec, err := e.embedder.EmbedSingle(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("embed query: %w", err)
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	var allResults []Result
	for docID, doc := range e.docs {
		for _, c := range doc.Chunks {
			score := cosineSimilarity(queryVec, c.Vec)
			if score < e.cfg.MinScore {
				continue
			}
			allResults = append(allResults, Result{
				DocID:   docID,
				ChunkID: c.Idx,
				Title:   c.Title,
				Score:   score,
				Offset:  c.Offset,
				Length:  c.Length,
			})
		}
	}

	slices.SortFunc(allResults, func(a, b Result) int {
		if a.Score > b.Score {
			return -1
		}
		if a.Score < b.Score {
			return 1
		}
		return 0
	})

	if len(allResults) > topK {
		allResults = allResults[:topK]
	}

	return positionAwareOrder(allResults), nil
}

func positionAwareOrder(results []Result) []Result {
	if len(results) <= 2 {
		return results
	}

	ordered := make([]Result, 0, len(results))

	top2 := results[:2]
	middle := results[2:]

	ordered = append(ordered, top2...)
	ordered = append(ordered, middle...)

	return ordered
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float32
	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dot / float32(math.Sqrt(float64(normA*normB)))
}

func (e *Engine) GetDocument(docID string) (*Document, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	doc, ok := e.docs[docID]
	return doc, ok
}

func (e *Engine) ListDocuments() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	ids := make([]string, 0, len(e.docs))
	for id := range e.docs {
		ids = append(ids, id)
	}
	return ids
}

func (e *Engine) Reload(dir string) error {
	newEngine := NewEngine(e.cfg, e.embedder)
	if err := newEngine.LoadFromDisk(dir); err != nil {
		return err
	}

	e.mu.Lock()
	e.docs = newEngine.docs
	e.mu.Unlock()

	return nil
}
