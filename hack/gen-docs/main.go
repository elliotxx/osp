package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/elliotxx/osp/pkg/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	// Parse command line flags
	outputDir := flag.String("output", "docs/reference", "Output directory for documentation")
	flag.Parse()

	// Get the root command
	rootCmd := cmd.GetRootCmd()

	// Create docs directory if not exists
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create docs directory: %v", err)
	}

	// Clean up old docs
	if err := cleanDir(*outputDir); err != nil {
		log.Fatalf("Failed to clean docs directory: %v", err)
	}

	// Generate markdown documentation
	if err := doc.GenMarkdownTree(rootCmd, *outputDir); err != nil {
		log.Fatalf("Failed to generate markdown docs: %v", err)
	}
}

// cleanDir removes all files in the specified directory
func cleanDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.RemoveAll(filepath.Join(dir, file.Name())); err != nil {
			return err
		}
	}

	return nil
}
