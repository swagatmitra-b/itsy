package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:  "itsy [url]",
		Args: cobra.ExactArgs(1),
		RunE: RootFunc,
	}

	rootCmd.Flags().StringVarP(&tags, "tags", "t", "", "Comma-separated list of tags to extract")
	rootCmd.Flags().BoolVarP(&output, "out", "o", false, "Save output to a text file")
	rootCmd.Flags().StringVarP(&css, "css", "c", "", "Comma-separated list of CSS selectors")
	rootCmd.Flags().IntVarP(&depth, "depth", "d", 0, "Crawl depth")
	rootCmd.Flags().StringVarP(&wordsearch, "wordsearch", "w", "", "Comma or semi-colon separated list of tags")
	rootCmd.Flags().BoolVarP(&internalOnly, "internal", "i", false, "Crawl only internal links (same domain and subdomains)")
	rootCmd.Flags().BoolVarP(&sitetree, "sitetree", "s", false, "Print the site-tree")

	imageCmd := &cobra.Command{
		Use:  "image [path]",
		Args: cobra.ExactArgs(1),
		RunE: Ascii,
	}

	imageCmd.Flags().IntVar(&width, "width", 100, "Width of the ASCII output")
	imageCmd.Flags().IntVar(&height, "height", 0, "Height of ASCII output")
	imageCmd.Flags().BoolVarP(&invert, "invert", "i", false, "Invert brightness (dark on light)")
	imageCmd.Flags().BoolVarP(&color, "color", "c", false, "Enable ANSI color output")
	imageCmd.Flags().StringVarP(&outFile, "out", "o", "", "Output file (leave empty for stdout)")

	rootCmd.AddCommand(imageCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
