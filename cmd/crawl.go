package main

import (
	"fmt"
	"net/url"
	"os"
	"scrawl/utils"

	"github.com/PuerkitoBio/goquery"
)

func Crawl(resource string, depth int) (int, int) {
	success, failed := 0, 0
	visited := make(map[string]bool)
	// tree := make(map[string][]string)

	type Page struct {
		url   string
		depth int
	}

	q := []Page{{url: resource, depth: 0}}

	for len(q) > 0 {
		page := q[0]
		q = q[1:]

		if visited[page.url] || page.depth > depth {
			continue
		}
		visited[page.url] = true

		fmt.Printf("\n=== Crawling [%s] at depth %d ===\n\n", page.url, page.depth)

		doc, err := processURL(page.url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to crawl %s: %v\n", page.url, err)
			failed += 1
			continue
		}

		utils.Wait()

		if doc == nil {
			continue
		}

		success += 1

		base, _ := url.Parse(page.url)
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if !exists {
				return
			}
			u, err := base.Parse(href)
			if err != nil {
				fmt.Println("Error in building link")
				return
			}
			absolute := u.String()
			if !visited[absolute] {
				q = append(q, Page{url: absolute, depth: page.depth + 1})
				// tree[page.url] = append(tree[page.url], absolute)
			}
		})
	}

	// if sitetree {
	// 	fmt.Println("\n=== Site Tree ===")
	// printTree()
	// }
	return success, failed
}
