package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type FileConfig struct {
	Provider string `json:"provider,omitempty"`
	Endpoint string `json:"endpoint"`
	APIKey   string `json:"api_key,omitempty"`
	Model    string `json:"model,omitempty"`
	DocsPath string `json:"docs_path"`
}

var skillFolderNames = []string{
	".opencode/skills/rag-doc-index",
	".claude/skills/rag-doc-index",
	"skills/rag-doc-index",
}

func findSkillFolder() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	dir := cwd
	for {
		for _, skillPath := range skillFolderNames {
			fullPath := filepath.Join(dir, skillPath)
			if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
				return fullPath
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

func getSkillFolder() string {
	if folder := findSkillFolder(); folder != "" {
		return folder
	}

	if envFolder := os.Getenv("RAG_SKILL_FOLDER"); envFolder != "" {
		return envFolder
	}

	cwd, _ := os.Getwd()
	for _, skillPath := range skillFolderNames {
		fullPath := filepath.Join(cwd, skillPath)
		os.MkdirAll(fullPath, 0755)
		return fullPath
	}

	return cwd
}

func getConfigPath() string {
	return filepath.Join(getSkillFolder(), "rag-doc-index.config.json")
}

func getEmbeddingsPath() string {
	return filepath.Join(getSkillFolder(), "embeddings")
}

func loadFileConfig() FileConfig {
	cfg := FileConfig{
		Provider: "",
		Endpoint: "",
		APIKey:   "",
		Model:    "",
		DocsPath: "",
	}

	path := getConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	json.Unmarshal(data, &cfg)
	return cfg
}

func saveFileConfig(cfg FileConfig) error {
	path := getConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func configValue(fileVal, envVal, defaultVal string) string {
	if fileVal != "" {
		return fileVal + " (from config file)"
	}
	if envVal != "" {
		return envVal + " (from env)"
	}
	return defaultVal + " (default)"
}

func maskAPIKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure settings",
	Long: `Configure embedding provider, endpoint, and documents path.

Settings are saved to rag-doc-index.config.json in the skill folder
(.opencode/skills/rag-doc-index or .claude/skills/rag-doc-index).

Supported providers:
  - tei: Text Embeddings Inference (HuggingFace)
  - omlx: oMLX (OpenAI-compatible, Apple Silicon optimized)
  - openai: OpenAI API`,
	Run: func(cmd *cobra.Command, args []string) {
		fileCfg := loadFileConfig()
		show, _ := cmd.Flags().GetBool("show")
		provider, _ := cmd.Flags().GetString("provider")
		tei, _ := cmd.Flags().GetString("tei")
		apiKey, _ := cmd.Flags().GetString("api-key")
		model, _ := cmd.Flags().GetString("model")
		docs, _ := cmd.Flags().GetString("docs")
		skillFolder, _ := cmd.Flags().GetString("skill-folder")
		changed := false

		if skillFolder != "" {
			os.Setenv("RAG_SKILL_FOLDER", skillFolder)
		}

		if provider != "" {
			fileCfg.Provider = provider
			changed = true
		}
		if tei != "" {
			fileCfg.Endpoint = tei
			changed = true
		}
		if apiKey != "" {
			fileCfg.APIKey = apiKey
			changed = true
		}
		if model != "" {
			fileCfg.Model = model
			changed = true
		}
		if docs != "" {
			fileCfg.DocsPath = docs
			changed = true
		}

		if show || !changed {
			fmt.Println("Current Configuration:")
			fmt.Printf("  Skill folder: %s\n", getSkillFolder())
			fmt.Printf("  Config file: %s\n", getConfigPath())
			fmt.Printf("  Embeddings: %s\n", getEmbeddingsPath())
			fmt.Println()

			providerVal := fileCfg.Provider
			if providerVal == "" {
				providerVal = os.Getenv("EMBEDDING_PROVIDER")
			}
			if providerVal == "" {
				providerVal = "tei"
			}
			fmt.Printf("  Provider: %s\n", providerVal)

			fmt.Printf("  Endpoint: %s\n", configValue(fileCfg.Endpoint, os.Getenv("TEI_ENDPOINT"), "http://host.docker.internal:8080"))

			apiKeyVal := fileCfg.APIKey
			if apiKeyVal == "" {
				apiKeyVal = os.Getenv("API_KEY")
			}
			fmt.Printf("  API Key: %s\n", maskAPIKey(apiKeyVal))

			modelVal := fileCfg.Model
			if modelVal == "" {
				modelVal = os.Getenv("EMBEDDING_MODEL")
			}
			if modelVal == "" {
				modelVal = "BAAI/bge-small-en-v1.5"
			}
			fmt.Printf("  Model: %s\n", modelVal)

			fmt.Printf("  Docs Path: %s\n", configValue(fileCfg.DocsPath, os.Getenv("DOCS_PATH"), "./docs"))
			return
		}

		if err := saveFileConfig(fileCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Configuration saved:")
		fmt.Printf("  Skill folder: %s\n", getSkillFolder())
		fmt.Printf("  Config file: %s\n", getConfigPath())
		fmt.Printf("  Provider: %s\n", fileCfg.Provider)
		fmt.Printf("  Endpoint: %s\n", fileCfg.Endpoint)
		fmt.Printf("  API Key: %s\n", maskAPIKey(fileCfg.APIKey))
		fmt.Printf("  Model: %s\n", fileCfg.Model)
		fmt.Printf("  Docs Path: %s\n", fileCfg.DocsPath)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().String("provider", "", "Embedding provider (tei, omlx, openai)")
	configCmd.Flags().String("tei", "", "Embedding endpoint URL")
	configCmd.Flags().String("api-key", "", "API key for authentication")
	configCmd.Flags().String("model", "", "Embedding model name")
	configCmd.Flags().String("docs", "", "Documents root path")
	configCmd.Flags().Bool("show", false, "Show current configuration")
	configCmd.Flags().String("skill-folder", "", "Override skill folder location")
}

func resolveDocsPath(docsPath string) string {
	if docsPath == "" {
		return "./docs"
	}

	if filepath.IsAbs(docsPath) {
		return docsPath
	}

	cwd, _ := os.Getwd()
	return filepath.Join(cwd, docsPath)
}

func getRelativePath(basePath, fullPath string) string {
	rel, err := filepath.Rel(basePath, fullPath)
	if err != nil {
		return fullPath
	}
	return strings.TrimPrefix(rel, "./")
}

func applyFileConfigToEnv(cfg FileConfig) {
	if cfg.Provider != "" {
		os.Setenv("EMBEDDING_PROVIDER", cfg.Provider)
	}
	if cfg.Endpoint != "" {
		os.Setenv("EMBEDDING_ENDPOINT", cfg.Endpoint)
	}
	if cfg.APIKey != "" {
		os.Setenv("API_KEY", cfg.APIKey)
	}
	if cfg.Model != "" {
		os.Setenv("EMBEDDING_MODEL", cfg.Model)
	}
	if cfg.DocsPath != "" {
		os.Setenv("DOCS_PATH", cfg.DocsPath)
	}
}
