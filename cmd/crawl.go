package main

import (
	"fmt"
	"os"
	"scrawl/utils"

	"github.com/PuerkitoBio/goquery"
)

func Crawl(resource string, depth int) (int, int) {
	success, failed := 0, 0
	visited := make(map[string]bool)
	tree := make(map[string][]string)

	baseURL, err := utils.NormalizeURL(resource)
	if err != nil {
		utils.ExitWithError("Failed to normalize URL", err)
		return 0, 0
	}
	resource = baseURL.String()
	baseDomain := baseURL.Hostname()

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

		base, _ := utils.NormalizeURL(page.url)
		seenChildren := make(map[string]bool)

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

			normalized, err := utils.NormalizeURL(u.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to normalize URL: %v\n", err)
				return
			}

			absolute := normalized.String()

			if internalOnly {
				if u.Hostname() != baseDomain && !utils.HasSubdomain(u.Hostname(), baseDomain) {
					return
				}
			}

			if !seenChildren[absolute] {
				seenChildren[absolute] = true
				tree[page.url] = append(tree[page.url], absolute)
			}

			if !visited[absolute] {
				q = append(q, Page{url: absolute, depth: page.depth + 1})
			}
		})
	}

	if sitetree {
		fmt.Println(tree)
		printTree(tree, map[string]bool{}, resource, "", 0, depth)
	}
	return success, failed
}

func printTree(tree map[string][]string, visited map[string]bool, node, indent string, currDepth, depth int) {

	if currDepth > depth {
		return
	}
	visited[node] = true
	fmt.Println(indent + node)

	for i, child := range tree[node] {
		connector := "├──"
		if i == len(tree[node])-1 {
			connector = "└──"
		}

		if !visited[child] {
			printTree(tree, visited, child, indent+connector, currDepth+1, depth)
		}
	}
}
