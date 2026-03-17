package embed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vrealzhou/doc-index/internal/config"
)

type Client struct {
	cfg    config.Config
	client *http.Client
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		cfg: cfg,
		client: &http.Client{
			Timeout: cfg.RequestTimeout,
		},
	}
}

type TEIRequest struct {
	Inputs    []string `json:"inputs"`
	Normalize bool     `json:"normalize,omitempty"`
}

type OpenAIEmbeddingRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

type OpenAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
}

func (c *Client) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	switch c.cfg.Provider {
	case config.ProviderTEI:
		return c.embedTEI(ctx, texts)
	case config.ProviderOMLX, config.ProviderOpenAI:
		return c.embedOpenAI(ctx, texts)
	default:
		return c.embedTEI(ctx, texts)
	}
}

func (c *Client) embedTEI(ctx context.Context, texts []string) ([][]float32, error) {
	reqBody, _ := json.Marshal(TEIRequest{
		Inputs:    texts,
		Normalize: true,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", c.cfg.Endpoint+"/embed", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tei request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tei error %d: %s", resp.StatusCode, string(body))
	}

	var vectors [][]float32
	if err := json.NewDecoder(resp.Body).Decode(&vectors); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(vectors) != len(texts) {
		return nil, fmt.Errorf("vector count mismatch: expected %d, got %d", len(texts), len(vectors))
	}

	return vectors, nil
}

func (c *Client) embedOpenAI(ctx context.Context, texts []string) ([][]float32, error) {
	model := c.cfg.Model
	if model == "" {
		model = "bge-small-en-v1.5"
	}

	reqBody, _ := json.Marshal(OpenAIEmbeddingRequest{
		Input: texts,
		Model: model,
	})

	endpoint := c.cfg.Endpoint + "/v1/embeddings"

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedding request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding error %d: %s", resp.StatusCode, string(body))
	}

	var embedResp OpenAIEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	vectors := make([][]float32, len(texts))
	for _, data := range embedResp.Data {
		if data.Index < len(vectors) {
			vectors[data.Index] = data.Embedding
		}
	}

	if len(vectors) != len(texts) {
		return nil, fmt.Errorf("vector count mismatch: expected %d, got %d", len(texts), len(vectors))
	}

	return vectors, nil
}

func (c *Client) EmbedSingle(ctx context.Context, text string) ([]float32, error) {
	vecs, err := c.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}
	return vecs[0], nil
}

func (c *Client) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	switch c.cfg.Provider {
	case config.ProviderTEI:
		return c.healthCheckTEI(ctx)
	case config.ProviderOMLX, config.ProviderOpenAI:
		return c.healthCheckOpenAI(ctx)
	default:
		return c.healthCheckTEI(ctx)
	}
}

func (c *Client) healthCheckTEI(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.cfg.Endpoint+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("tei unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("tei health check failed: status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) healthCheckOpenAI(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.cfg.Endpoint+"/v1/models", nil)
	if err != nil {
		return err
	}

	if c.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("endpoint unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("health check failed: status %d", resp.StatusCode)
	}

	return nil
}
