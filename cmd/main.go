package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

func main() {
    rootCmd := &cobra.Command{
        Use:  "itsy [command]",
        Args:  cobra.ExactArgs(1),
		RunE: RootFunc,
	}

	rootCmd.Flags().StringVarP(&tags, "tags", "t", "", "Comma-separated list of tags to extract")
	rootCmd.Flags().BoolVarP(&output, "output", "o", false, "Save output to a file")
	rootCmd.Flags().StringVarP(&css, "css", "c", "", "Comma-separated list of CSS selectors")
	rootCmd.Flags().IntVarP(&depth, "depth", "d", 0, "Crawl depth")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
