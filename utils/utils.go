package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func generateFilename(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Sprintf("invalidurl_%x.txt", sha1.Sum([]byte(rawURL)))
	}

	fullPath := parsed.Host + parsed.Path

	safePath := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, fullPath)

	if len(safePath) > 100 {
		safePath = safePath[:100]
	}

	h := sha1.New()
	h.Write([]byte(rawURL))
	hash := hex.EncodeToString(h.Sum(nil))[:11]

	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s_%s.txt", safePath, hash, timestamp)
}

func ExitWithError(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(1)
}

func OutputPage(outstring, resource string) {
	dir := "pages"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		ExitWithError("Failed to create output directory", err)
	}

	filename := generateFilename(resource)
	fullPath := fmt.Sprintf("%s/%s", dir, filename)

	err = os.WriteFile(fullPath, []byte(outstring), 0644)
	if err != nil {
		ExitWithError("Failed to write to file", err)
	}
	fmt.Printf("Output saved to %s\n", fullPath)
}

func Wait() {
	minDelay := 700
	maxDelay := 1500
	delay := rand.Intn(maxDelay-minDelay) + minDelay
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func ParseSearchTerms(input string) (terms []string, matchAll bool) {
	input = strings.TrimSpace(input)
	matchAll = strings.Contains(input, "&")
	var split []string

	if matchAll {
		split = strings.Split(input, "&")
	} else {
		split = strings.Split(input, ",")
	}

	for _, term := range split {
		trimmed := strings.ToLower(strings.TrimSpace(term))
		if trimmed != "" {
			terms = append(terms, trimmed)
		}
	}
	return
}

func HasSubdomain(host, base string) bool {
	return len(host) > len(base) && host[len(host)-len(base)-1:] == "."+base
}

func NormalizeURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	q := u.Query()
	for _, junk := range []string{"utm_source", "utm_medium", "utm_campaign", "ref"} {
		q.Del(junk)
	}
	keys := make([]string, 0, len(q))
	for k := range q {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sorted := url.Values{}
	for _, k := range keys {
		sorted[k] = q[k]
	}
	u.RawQuery = sorted.Encode()

	u.Fragment = ""

	if u.Path == "" {
		u.Path = "/"
	}

	return u, nil
}

func ContainsWord(body, word string) bool {
	bodyLower := strings.ToLower(body)
	wordLower := regexp.QuoteMeta(strings.ToLower(word))

	pattern := `(?i)(^|[\s.,;:!?"'()\[\]{}])` + wordLower + `($|[\s.,;:!?"'()\[\]{}])`

	re := regexp.MustCompile(pattern)
	return re.MatchString(bodyLower)
}
