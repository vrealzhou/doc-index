package indexer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vrealzhou/doc-index/internal/config"
	"github.com/vrealzhou/doc-index/internal/embed"
)

func TestComputeHash(t *testing.T) {
	content := []byte("test content")
	hash := computeHash(content)

	if len(hash) != 16 {
		t.Errorf("expected hash length 16, got %d", len(hash))
	}

	sameHash := computeHash(content)
	if hash != sameHash {
		t.Error("hash should be deterministic")
	}

	differentContent := []byte("different content")
	differentHash := computeHash(differentContent)
	if hash == differentHash {
		t.Error("different content should produce different hash")
	}
}

func TestScanEmpty(t *testing.T) {
	docsDir, err := os.MkdirTemp("", "docs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(docsDir)

	embedDir, err := os.MkdirTemp("", "embeddings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(embedDir)

	cfg := config.Config{
		DocsPath:       docsDir,
		EmbeddingsPath: embedDir,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := embed.NewClient(cfg)
	idx := New(cfg, client)

	report, err := idx.Scan()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.Total != 0 {
		t.Errorf("expected 0 total documents, got %d", report.Total)
	}
}

func TestScanSingleDocument(t *testing.T) {
	docsDir, err := os.MkdirTemp("", "docs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(docsDir)

	docContent := `# Test Document

This is a test document with some content.

## Section One

More content here.`

	err = os.WriteFile(filepath.Join(docsDir, "test.md"), []byte(docContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	embedDir, err := os.MkdirTemp("", "embeddings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(embedDir)

	cfg := config.Config{
		DocsPath:       docsDir,
		EmbeddingsPath: embedDir,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := embed.NewClient(cfg)
	idx := New(cfg, client)

	report, err := idx.Scan()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.Total != 1 {
		t.Errorf("expected 1 total document, got %d", report.Total)
	}

	if report.Missing != 1 {
		t.Errorf("expected 1 missing document, got %d", report.Missing)
	}

	docsBase := filepath.Base(docsDir)
	docPath := docsBase + "/test.md"
	if _, ok := report.Details[docPath]; !ok {
		t.Errorf("expected document %s in details", docPath)
	}
}

func TestScanSubfolder(t *testing.T) {
	docsDir, err := os.MkdirTemp("", "docs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(docsDir)

	subDir := filepath.Join(docsDir, "subfolder")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(subDir, "nested.md"), []byte("# Nested"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(docsDir, "root.md"), []byte("# Root"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	embedDir, err := os.MkdirTemp("", "embeddings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(embedDir)

	cfg := config.Config{
		DocsPath:       docsDir,
		EmbeddingsPath: embedDir,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := embed.NewClient(cfg)
	idx := New(cfg, client)

	report, err := idx.Scan()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.Total != 2 {
		t.Errorf("expected 2 total documents, got %d", report.Total)
	}

	docsBase := filepath.Base(docsDir)
	if _, ok := report.Details[docsBase+"/root.md"]; !ok {
		t.Errorf("expected root.md in details")
	}
	if _, ok := report.Details[docsBase+"/subfolder/nested.md"]; !ok {
		t.Errorf("expected subfolder/nested.md in details")
	}
}

func TestScanStaleDocument(t *testing.T) {
	docsDir, err := os.MkdirTemp("", "docs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(docsDir)

	docContent := "# Test Document"
	err = os.WriteFile(filepath.Join(docsDir, "test.md"), []byte(docContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	embedDir, err := os.MkdirTemp("", "embeddings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(embedDir)

	docsBase := filepath.Base(docsDir)
	embedSubDir := filepath.Join(embedDir, docsBase)
	os.MkdirAll(embedSubDir, 0755)

	meta := map[string]interface{}{
		"meta":   true,
		"v":      1,
		"doc":    docsBase + "/test.md",
		"hash":   "differenthash123456",
		"chunks": 1,
	}
	metaJSON, _ := json.Marshal(meta)

	chunk := map[string]interface{}{
		"idx":    0,
		"title":  "Test",
		"offset": 0,
		"length": 100,
		"vec":    []float32{0.1, 0.2},
	}
	chunkJSON, _ := json.Marshal(chunk)

	jsonlContent := string(metaJSON) + "\n" + string(chunkJSON) + "\n"
	err = os.WriteFile(filepath.Join(embedSubDir, "test.md.jsonl"), []byte(jsonlContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg := config.Config{
		DocsPath:       docsDir,
		EmbeddingsPath: embedDir,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := embed.NewClient(cfg)
	idx := New(cfg, client)

	report, err := idx.Scan()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.Stale != 1 {
		t.Errorf("expected 1 stale document, got %d", report.Stale)
	}
}

func TestScanOrphanEmbedding(t *testing.T) {
	docsDir, err := os.MkdirTemp("", "docs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(docsDir)

	embedDir, err := os.MkdirTemp("", "embeddings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(embedDir)

	docsBase := filepath.Base(docsDir)
	embedSubDir := filepath.Join(embedDir, docsBase)
	os.MkdirAll(embedSubDir, 0755)

	meta := map[string]interface{}{
		"meta":   true,
		"v":      1,
		"doc":    docsBase + "/orphan.md",
		"hash":   "somehash12345678",
		"chunks": 1,
	}
	metaJSON, _ := json.Marshal(meta)

	chunk := map[string]interface{}{
		"idx":    0,
		"title":  "Orphan",
		"offset": 0,
		"length": 100,
		"vec":    []float32{0.1, 0.2},
	}
	chunkJSON, _ := json.Marshal(chunk)

	jsonlContent := string(metaJSON) + "\n" + string(chunkJSON) + "\n"
	err = os.WriteFile(filepath.Join(embedSubDir, "orphan.md.jsonl"), []byte(jsonlContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg := config.Config{
		DocsPath:       docsDir,
		EmbeddingsPath: embedDir,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := embed.NewClient(cfg)
	idx := New(cfg, client)

	report, err := idx.Scan()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.Orphan != 1 {
		t.Errorf("expected 1 orphan embedding, got %d", report.Orphan)
	}
}

func TestReindex(t *testing.T) {
	docsDir, err := os.MkdirTemp("", "docs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(docsDir)

	docContent := `# Test

Some test content here.`
	err = os.WriteFile(filepath.Join(docsDir, "test.md"), []byte(docContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	embedDir, err := os.MkdirTemp("", "embeddings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(embedDir)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vectors := [][]float32{{0.1, 0.2, 0.3}}
		json.NewEncoder(w).Encode(vectors)
	}))
	defer server.Close()

	cfg := config.Config{
		DocsPath:       docsDir,
		EmbeddingsPath: embedDir,
		Endpoint:       server.URL,
		Provider:       config.ProviderTEI,
		Model:          "test-model",
		VectorDim:      3,
		MaxChunkSize:   1000,
		OverlapSize:    100,
		RequestTimeout: 5 * time.Second,
	}

	client := embed.NewClient(cfg)
	idx := New(cfg, client)

	docsBase := filepath.Base(docsDir)
	docPath := docsBase + "/test.md"

	err = idx.Reindex(context.Background(), []string{docPath})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	embedFile := filepath.Join(embedDir, docPath+".jsonl")
	if _, err := os.Stat(embedFile); os.IsNotExist(err) {
		t.Errorf("expected embedding file at %s", embedFile)
	}
}

func TestAutoReindex(t *testing.T) {
	docsDir, err := os.MkdirTemp("", "docs")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(docsDir)

	err = os.WriteFile(filepath.Join(docsDir, "test.md"), []byte("# Test\n\nContent"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	embedDir, err := os.MkdirTemp("", "embeddings")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(embedDir)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vectors := [][]float32{{0.1, 0.2}}
		json.NewEncoder(w).Encode(vectors)
	}))
	defer server.Close()

	cfg := config.Config{
		DocsPath:       docsDir,
		EmbeddingsPath: embedDir,
		Endpoint:       server.URL,
		Provider:       config.ProviderTEI,
		MaxChunkSize:   1000,
		OverlapSize:    100,
		RequestTimeout: 5 * time.Second,
	}

	client := embed.NewClient(cfg)
	idx := New(cfg, client)

	report, err := idx.AutoReindex(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if report.Current != 1 {
		t.Errorf("expected 1 current document, got %d", report.Current)
	}
}

func TestReadStoredHash(t *testing.T) {
	content := `{"meta":true,"v":1,"doc":"test.md","hash":"abc123def456","chunks":1}
{"idx":0,"title":"Test","offset":0,"length":100}`

	tmpFile, err := os.CreateTemp("", "test*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(content)
	tmpFile.Close()

	hash, err := readStoredHash(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash != "abc123def456" {
		t.Errorf("expected hash abc123def456, got %s", hash)
	}
}

func TestReadStoredHashInvalid(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("not valid json")
	tmpFile.Close()

	_, err = readStoredHash(tmpFile.Name())
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestReadStoredHashEmpty(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Close()

	_, err = readStoredHash(tmpFile.Name())
	if err == nil {
		t.Error("expected error for empty file")
	}
}
