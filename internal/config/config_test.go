package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg := Load()

	if cfg.Endpoint != "http://host.docker.internal:8080" {
		t.Errorf("expected default endpoint, got %s", cfg.Endpoint)
	}
	if cfg.Model != "BAAI/bge-small-en-v1.5" {
		t.Errorf("expected default model, got %s", cfg.Model)
	}
	if cfg.VectorDim != 384 {
		t.Errorf("expected default vector dim 384, got %d", cfg.VectorDim)
	}
	if cfg.ServerPort != 8081 {
		t.Errorf("expected default server port 8081, got %d", cfg.ServerPort)
	}
	if cfg.DefaultTopK != 5 {
		t.Errorf("expected default top_k 5, got %d", cfg.DefaultTopK)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("EMBEDDING_ENDPOINT", "http://custom:9000")
	os.Setenv("EMBEDDING_DIM", "512")
	os.Setenv("SERVER_PORT", "9090")
	defer func() {
		os.Unsetenv("EMBEDDING_ENDPOINT")
		os.Unsetenv("EMBEDDING_DIM")
		os.Unsetenv("SERVER_PORT")
	}()

	cfg := Load()

	if cfg.Endpoint != "http://custom:9000" {
		t.Errorf("expected custom endpoint, got %s", cfg.Endpoint)
	}
	if cfg.VectorDim != 512 {
		t.Errorf("expected vector dim 512, got %d", cfg.VectorDim)
	}
	if cfg.ServerPort != 9090 {
		t.Errorf("expected server port 9090, got %d", cfg.ServerPort)
	}
}

func TestEnvString(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback string
		want     string
	}{
		{"env set", "TEST_KEY", "value", "default", "value"},
		{"env unset", "NONEXISTENT_KEY", "", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			got := envString(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("envString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvInt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback int
		want     int
	}{
		{"valid int", "TEST_INT", "42", 0, 42},
		{"invalid int", "TEST_INT", "notanumber", 10, 10},
		{"unset", "NONEXISTENT_INT", "", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			got := envInt(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("envInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvFloat(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback float32
		want     float32
	}{
		{"valid float", "TEST_FLOAT", "0.5", 0.0, 0.5},
		{"invalid float", "TEST_FLOAT", "notanumber", 0.3, 0.3},
		{"unset", "NONEXISTENT_FLOAT", "", 0.7, 0.7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			got := envFloat(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("envFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
