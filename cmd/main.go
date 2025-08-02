package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "itsy [command]",
		Args: cobra.ExactArgs(1),
		RunE: RootFunc,
	}

	rootCmd.Flags().StringVarP(&tags, "tags", "t", "", "Comma-separated list of tags to extract")
	rootCmd.Flags().BoolVarP(&output, "output", "o", false, "Save output to a file")
	rootCmd.Flags().StringVarP(&css, "css", "c", "", "Comma-separated list of CSS selectors")
	rootCmd.Flags().IntVarP(&depth, "depth", "d", 0, "Crawl depth")
	rootCmd.Flags().StringVarP(&wordsearch, "wordsearch", "w", "", "Comma or semi-colon separated list of tags")
	rootCmd.Flags().BoolVarP(&internalOnly, "internal", "i", false, "Crawl only internal links (same domain and subdomains)")
	rootCmd.Flags().BoolVarP(&sitetree, "sitetree", "s", false, "Print the site-tree")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
