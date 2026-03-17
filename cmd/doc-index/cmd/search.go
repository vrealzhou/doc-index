package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/vrealzhou/doc-index/internal/search"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search documents with semantic query",
	Long: `Search for relevant documentation sections using natural language.

Results include document ID, section title, position, and similarity score.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _, _, engine := loadComponents()

		query := args[0]
		for i := 1; i < len(args); i++ {
			query += " " + args[i]
		}

		topK, _ := cmd.Flags().GetInt("top-k")
		minScore, _ := cmd.Flags().GetFloat32("min-score")
		outputJSON, _ := cmd.Flags().GetBool("json")

		if topK <= 0 {
			topK = cfg.DefaultTopK
		}
		if minScore == 0 {
			minScore = cfg.MinScore
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		results, err := engine.Search(ctx, query, topK)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching: %v\n", err)
			os.Exit(1)
		}

		var filtered []search.Result
		for _, r := range results {
			if r.Score >= minScore {
				filtered = append(filtered, r)
			}
		}

		if outputJSON {
			output := map[string]interface{}{
				"query":   query,
				"total":   len(filtered),
				"results": filtered,
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			enc.Encode(output)
			return
		}

		if len(filtered) == 0 {
			fmt.Printf("No results found for: %s\n", query)
			fmt.Printf("Try lowering --min-score (current: %.2f)\n", minScore)
			return
		}

		fmt.Printf("Found %d results for: %s\n\n", len(filtered), query)

		for i, r := range filtered {
			fmt.Printf("--- Result %d (score: %.3f) ---\n", i+1, r.Score)
			fmt.Printf("Document: %s\n", r.DocID)
			fmt.Printf("Section: %s\n", r.Title)
			fmt.Printf("Position: offset=%d, length=%d\n", r.Offset, r.Length)
			if r.Preview != "" {
				fmt.Printf("Preview: %s\n", r.Preview)
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntP("top-k", "k", 0, "Number of results (default: 5)")
	searchCmd.Flags().Float32P("min-score", "m", 0, "Minimum similarity score (default: 0.3)")
	searchCmd.Flags().BoolP("json", "j", false, "Output as JSON")
}
