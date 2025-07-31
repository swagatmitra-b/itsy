package main

import (
	"fmt"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	"os"
)

func Crawl(resource string, depth int) (int, int) {
	success, failed := 0, 0

	type Page struct {
		url   string
		depth int
	}

	visited := make(map[string]bool)
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
			}
		})
	}
	return success, failed
}
