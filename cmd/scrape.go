package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"scrawl/utils"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var (
	tags         string
	output       bool
	css          string
	depth        int
	internalOnly bool
	wordsearch   string
	sitetree     bool
)

func RootFunc(cmd *cobra.Command, args []string) error {
	resource := args[0]

	if depth > 0 {
		success, failed := Crawl(resource, depth)
		fmt.Println()
		fmt.Printf("Success: %d, Failed: %d\n", success, failed)
		fmt.Printf("Total pages: %d", success+failed)
	} else {
		processURL(resource)
	}
	return nil
}

func processURL(resource string) (*goquery.Document, error) {
	resp, err := http.Get(resource)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyBytes)

	if wordsearch != "" {
		terms, matchAll := utils.ParseSearchTerms(wordsearch)

		found := 0
		for _, word := range terms {
			if utils.ContainsWord(bodyString, word) {
				found += 1
			}
		}

		if (matchAll && found != len(terms)) || (!matchAll && found == 0) {
			fmt.Println("No matches")
			return nil, nil
		}
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyString))
	if err != nil {
		fmt.Println("Unable to parse DOM")
		return nil, err
	}
	base, _ := url.Parse(resource)

	var result strings.Builder

	if tags != "" {
		tagList := strings.Split(tags, ",")
		for _, tag := range tagList {
			tag = strings.TrimSpace(tag)
			if tag == "" {
				continue
			}
			appendMatches(&result, doc.Find(tag), fmt.Sprintf("<%s>", tag), base)
		}
	}

	if css != "" {
		selectorList := strings.Split(css, ",")
		for _, selector := range selectorList {
			selector = strings.TrimSpace(selector)
			if selector == "" {
				continue
			}
			appendMatches(&result, doc.Find(selector), fmt.Sprintf("Selector: %s", selector), base)
		}
	}

	if tags == "" && css == "" {
		if !output {
			fmt.Println(bodyString)
		} else {
			utils.OutputPage(bodyString, resource)
		}
		return doc, nil
	}

	if output {
		utils.OutputPage(result.String(), resource)
	} else {
		fmt.Println(result.String())
	}

	return doc, nil
}

func appendMatches(b *strings.Builder, selection *goquery.Selection, label string, base *url.URL) {
	fmt.Fprintf(b, "%s:\n\n", label)
	count := selection.Length()

	selection.Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text == "" {
			return
		}

		tagName := goquery.NodeName(s)
		if tagName == "a" {
			href, exists := s.Attr("href")
			if exists {
				absoluteURL := href
				if base != nil {
					u, err := base.Parse(href)
					if err == nil {
						absoluteURL = u.String()
					}
				}

				fmt.Fprintf(b, "- (%d/%d)\n%s\nLink: %s\n\n", i+1, count, text, absoluteURL)
				return
			}
		}

		fmt.Fprintf(b, "- (%d/%d)\n%s\n\n", i+1, count, text)
	})
}
