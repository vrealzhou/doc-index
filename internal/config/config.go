package config

import (
	"os"
	"strconv"
	"time"
)

type Provider string

const (
	ProviderTEI    Provider = "tei"
	ProviderOMLX   Provider = "omlx"
	ProviderOpenAI Provider = "openai"
)

type Config struct {
	Provider         Provider      `json:"provider"`
	Endpoint         string        `json:"endpoint"`
	APIKey           string        `json:"api_key,omitempty"`
	Model            string        `json:"model"`
	VectorDim        int           `json:"vector_dim"`
	BatchSize        int           `json:"batch_size"`
	RequestTimeout   time.Duration `json:"request_timeout_ms"`
	DocsPath         string        `json:"docs_path"`
	EmbeddingsPath   string        `json:"embeddings_path"`
	MaxChunkSize     int           `json:"max_chunk_size"`
	OverlapSize      int           `json:"overlap_size"`
	MaxContextTokens int           `json:"max_context_tokens"`
	WarnThreshold    int           `json:"warn_threshold"`
	HardLimit        int           `json:"hard_limit"`
	ServerPort       int           `json:"server_port"`
	DefaultTopK      int           `json:"default_top_k"`
	MinScore         float32       `json:"min_score"`
	PreviewLen       int           `json:"preview_len"`
}

func Load() Config {
	provider := Provider(envString("EMBEDDING_PROVIDER", "tei"))
	if provider != ProviderTEI && provider != ProviderOMLX && provider != ProviderOpenAI {
		provider = ProviderTEI
	}

	return Config{
		Provider:         provider,
		Endpoint:         envString("EMBEDDING_ENDPOINT", "http://host.docker.internal:8080"),
		APIKey:           envString("API_KEY", ""),
		Model:            envString("EMBEDDING_MODEL", "BAAI/bge-small-en-v1.5"),
		VectorDim:        envInt("EMBEDDING_DIM", 384),
		BatchSize:        envInt("BATCH_SIZE", 32),
		RequestTimeout:   time.Duration(envInt("REQUEST_TIMEOUT_MS", 30000)) * time.Millisecond,
		DocsPath:         envString("DOCS_PATH", "./docs"),
		EmbeddingsPath:   envString("EMBEDDINGS_PATH", "./embeddings"),
		MaxChunkSize:     envInt("MAX_CHUNK_SIZE", 1200),
		OverlapSize:      envInt("OVERLAP_SIZE", 100),
		MaxContextTokens: envInt("MAX_CONTEXT_TOKENS", 20000),
		WarnThreshold:    envInt("WARN_THRESHOLD", 16000),
		HardLimit:        envInt("HARD_LIMIT", 18000),
		ServerPort:       envInt("SERVER_PORT", 8081),
		DefaultTopK:      envInt("DEFAULT_TOP_K", 5),
		MinScore:         envFloat("MIN_SCORE", 0.3),
		PreviewLen:       envInt("PREVIEW_LEN", 150),
	}
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func envFloat(key string, fallback float32) float32 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(f)
		}
	}
	return fallback
}
