package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrealzhou/doc-index/internal/indexer"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show index status",
	Long: `Show current index status including document counts,
TEI connectivity, and per-document status.`,
	Run: func(cmd *cobra.Command, args []string) {
		skillFolder, _ := cmd.Flags().GetString("skill-folder")
		if skillFolder != "" {
			os.Setenv("RAG_SKILL_FOLDER", skillFolder)
		}

		cfg, embedder, idx, engine := loadComponents()

		report, err := idx.Scan()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning documents: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("=== Document Index Status ===")
		fmt.Printf("Skill folder: %s\n", getSkillFolder())
		fmt.Printf("Config file: %s\n", getConfigPath())
		fmt.Printf("Documents Path: %s\n", cfg.DocsPath)
		fmt.Printf("Embeddings Path: %s\n", cfg.EmbeddingsPath)
		fmt.Printf("Provider: %s\n", cfg.Provider)
		fmt.Printf("Endpoint: %s\n", cfg.Endpoint)
		fmt.Printf("Model: %s (%d dimensions)\n", cfg.Model, cfg.VectorDim)
		if cfg.APIKey != "" {
			fmt.Printf("API Key: %s\n", maskAPIKey(cfg.APIKey))
		}
		fmt.Println()

		teiStatus := "ok"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := embedder.HealthCheck(ctx); err != nil {
			teiStatus = "unavailable: " + err.Error()
		}
		fmt.Printf("TEI Status: %s\n", teiStatus)
		fmt.Println()

		fmt.Printf("Total Documents: %d\n", report.Total)
		fmt.Printf("  Current: %d\n", report.Current)
		fmt.Printf("  Stale: %d\n", report.Stale)
		fmt.Printf("  Missing: %d\n", report.Missing)
		fmt.Printf("  Orphan: %d\n", report.Orphan)
		fmt.Println()

		docsLoaded := engine.ListDocuments()
		fmt.Printf("Documents Loaded in Memory: %d\n", len(docsLoaded))

		if len(report.Details) > 0 {
			fmt.Println("\nDocument Details:")
			for name, status := range report.Details {
				statusStr := "current"
				switch status {
				case indexer.StatusStale:
					statusStr = "stale"
				case indexer.StatusMissing:
					statusStr = "missing"
				case indexer.StatusOrphan:
					statusStr = "orphan"
				}
				fmt.Printf("  %s: %s\n", name, statusStr)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().String("skill-folder", "", "Override skill folder location")
}
