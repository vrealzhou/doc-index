package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "doc-index",
	Short: "RAG document search CLI",
	Long: `doc-index provides semantic search over markdown documentation
using RAG (Retrieval-Augmented Generation) capabilities.

It indexes markdown documents and provides progressive disclosure
of design details to coding agents without blowing up context windows.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "Help for doc-index")
}
