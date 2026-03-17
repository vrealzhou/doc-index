package embed

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vrealzhou/doc-index/internal/config"
)

func TestEmbedBatchEmpty(t *testing.T) {
	cfg := config.Config{
		Provider:       config.ProviderTEI,
		Endpoint:       "http://localhost:8080",
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	result, err := client.EmbedBatch(context.Background(), []string{})
	if err != nil {
		t.Errorf("expected no error for empty input, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for empty input, got %v", result)
	}
}

func TestEmbedBatchTEI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/embed" {
			t.Errorf("expected /embed path, got %s", r.URL.Path)
		}

		var req TEIRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if len(req.Inputs) != 2 {
			t.Errorf("expected 2 inputs, got %d", len(req.Inputs))
		}

		vectors := [][]float32{
			{0.1, 0.2, 0.3},
			{0.4, 0.5, 0.6},
		}
		json.NewEncoder(w).Encode(vectors)
	}))
	defer server.Close()

	cfg := config.Config{
		Provider:       config.ProviderTEI,
		Endpoint:       server.URL,
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	result, err := client.EmbedBatch(context.Background(), []string{"hello", "world"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 vectors, got %d", len(result))
	}

	if len(result[0]) != 3 {
		t.Errorf("expected 3 dimensions, got %d", len(result[0]))
	}
}

func TestEmbedBatchTEIWithAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Bearer test-key, got %s", auth)
		}

		vectors := [][]float32{{0.1, 0.2}}
		json.NewEncoder(w).Encode(vectors)
	}))
	defer server.Close()

	cfg := config.Config{
		Provider:       config.ProviderTEI,
		Endpoint:       server.URL,
		APIKey:         "test-key",
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	_, err := client.EmbedBatch(context.Background(), []string{"test"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestEmbedBatchOpenAI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/embeddings" {
			t.Errorf("expected /v1/embeddings path, got %s", r.URL.Path)
		}

		var req OpenAIEmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}

		if req.Model != "test-model" {
			t.Errorf("expected model test-model, got %s", req.Model)
		}

		resp := OpenAIEmbeddingResponse{
			Data: []struct {
				Embedding []float32 `json:"embedding"`
				Index     int       `json:"index"`
			}{
				{Embedding: []float32{0.1, 0.2}, Index: 0},
				{Embedding: []float32{0.3, 0.4}, Index: 1},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := config.Config{
		Provider:       config.ProviderOMLX,
		Endpoint:       server.URL,
		Model:          "test-model",
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	result, err := client.EmbedBatch(context.Background(), []string{"a", "b"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 vectors, got %d", len(result))
	}
}

func TestEmbedBatchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	cfg := config.Config{
		Provider:       config.ProviderTEI,
		Endpoint:       server.URL,
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	_, err := client.EmbedBatch(context.Background(), []string{"test"})
	if err == nil {
		t.Error("expected error for 500 response")
	}
}

func TestEmbedSingle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vectors := [][]float32{{0.1, 0.2, 0.3}}
		json.NewEncoder(w).Encode(vectors)
	}))
	defer server.Close()

	cfg := config.Config{
		Provider:       config.ProviderTEI,
		Endpoint:       server.URL,
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	result, err := client.EmbedSingle(context.Background(), "test")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 dimensions, got %d", len(result))
	}
}

func TestHealthCheckTEI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("expected /health path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.Config{
		Provider:       config.ProviderTEI,
		Endpoint:       server.URL,
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	err := client.HealthCheck(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHealthCheckOMLX(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/models" {
			t.Errorf("expected /v1/models path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.Config{
		Provider:       config.ProviderOMLX,
		Endpoint:       server.URL,
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	err := client.HealthCheck(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHealthCheckWithAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Bearer test-key, got %s", auth)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	cfg := config.Config{
		Provider:       config.ProviderOMLX,
		Endpoint:       server.URL,
		APIKey:         "test-key",
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	err := client.HealthCheck(context.Background())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestHealthCheckFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	cfg := config.Config{
		Provider:       config.ProviderTEI,
		Endpoint:       server.URL,
		RequestTimeout: 5 * time.Second,
	}
	client := NewClient(cfg)

	err := client.HealthCheck(context.Background())
	if err == nil {
		t.Error("expected error for unhealthy endpoint")
	}
}
