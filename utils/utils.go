package utils

import (
	"fmt"
	"time"
	"strings"
	"os"
	"crypto/sha1"
	"encoding/hex"
)

func generateFilename(url string) string {
	
	h := sha1.New()
	h.Write([]byte(url))
	hash := hex.EncodeToString(h.Sum(nil))[:11] 

	safePrefix := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, url)

	if len(safePrefix) > 50 {
		safePrefix = safePrefix[:50]
	}

	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s_%s.txt", safePrefix, hash, timestamp)
}


func exitWithError(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(1)
}

func OutputPage(outstring, resource string) {
	dir := "pages"
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		exitWithError("Failed to create output directory", err)
	}

	filename := generateFilename(resource)
	fullPath := fmt.Sprintf("%s/%s", dir, filename)

	err = os.WriteFile(fullPath, []byte(outstring), 0644)
	if err != nil {
		exitWithError("Failed to write to file", err)
	}
	fmt.Printf("Output saved to %s\n", fullPath)
}