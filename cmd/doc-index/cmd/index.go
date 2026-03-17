package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrealzhou/doc-index/internal/config"
	"github.com/vrealzhou/doc-index/internal/embed"
	"github.com/vrealzhou/doc-index/internal/indexer"
	"github.com/vrealzhou/doc-index/internal/search"
)

func loadComponents() (config.Config, *embed.Client, *indexer.Indexer, *search.Engine) {
	fileCfg := loadFileConfig()
	applyFileConfigToEnv(fileCfg)

	docsPath := fileCfg.DocsPath
	if docsPath == "" {
		docsPath = os.Getenv("DOCS_PATH")
	}
	if docsPath == "" {
		docsPath = "./docs"
	}
	docsPath = resolveDocsPath(docsPath)

	embeddingsPath := getEmbeddingsPath()

	os.Setenv("DOCS_PATH", docsPath)
	os.Setenv("EMBEDDINGS_PATH", embeddingsPath)

	cfg := config.Load()
	embedder := embed.NewClient(cfg)
	idx := indexer.New(cfg, embedder)
	engine := search.NewEngine(cfg, embedder)

	if err := engine.LoadFromDisk(cfg.EmbeddingsPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load embeddings: %v\n", err)
	}

	return cfg, embedder, idx, engine
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index all markdown documents",
	Long: `Index all markdown documents in the configured docs path.

Documents are chunked, embedded via TEI, and stored as JSONL files
in the embeddings directory within the skill folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		skillFolder, _ := cmd.Flags().GetString("skill-folder")
		if skillFolder != "" {
			os.Setenv("RAG_SKILL_FOLDER", skillFolder)
		}

		cfg, _, idx, engine := loadComponents()
		force, _ := cmd.Flags().GetBool("force")

		report, err := idx.Scan()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning documents: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Skill folder: %s\n", getSkillFolder())
		fmt.Printf("Docs path: %s\n", cfg.DocsPath)
		fmt.Printf("Embeddings: %s\n", cfg.EmbeddingsPath)
		fmt.Println()
		fmt.Printf("Documents: %d total, %d current, %d stale, %d missing, %d orphan\n",
			report.Total, report.Current, report.Stale, report.Missing, report.Orphan)

		var toIndex []string
		if force {
			for name := range report.Details {
				toIndex = append(toIndex, name)
			}
		} else {
			for name, status := range report.Details {
				if status == indexer.StatusStale || status == indexer.StatusMissing {
					toIndex = append(toIndex, name)
				}
			}
		}

		for name, status := range report.Details {
			if status == indexer.StatusOrphan {
				fmt.Printf("Removing orphan: %s\n", name)
			}
		}

		if len(toIndex) == 0 {
			fmt.Println("All documents are up to date. Use --force to reindex.")
			return
		}

		fmt.Printf("Indexing %d documents...\n", len(toIndex))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		if force {
			if err := idx.Reindex(ctx, toIndex); err != nil {
				fmt.Fprintf(os.Stderr, "Error reindexing: %v\n", err)
				os.Exit(1)
			}
		} else {
			if _, err := idx.AutoReindex(ctx); err != nil {
				fmt.Fprintf(os.Stderr, "Error auto-reindexing: %v\n", err)
				os.Exit(1)
			}
		}

		if err := engine.Reload(cfg.EmbeddingsPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error reloading search index: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Indexing complete.")

		report, _ = idx.Scan()
		fmt.Printf("Final status: %d current, %d stale, %d missing, %d orphan\n",
			report.Current, report.Stale, report.Missing, report.Orphan)
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
	indexCmd.Flags().BoolP("force", "f", false, "Force reindex all documents")
	indexCmd.Flags().String("skill-folder", "", "Override skill folder location")
}
