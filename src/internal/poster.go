package internal

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// GeneratePoster generates a poster image using arxiv2poster for a given paper
// Returns the path to the generated poster, or empty string if generation failed
func GeneratePoster(paperLink string, outputDir string) string {
	// Extract arXiv ID from the paper link
	// HuggingFace paper links are like: https://huggingface.co/papers/2512.13043
	// We need to convert this to arXiv format

	// Check if this is an arXiv paper (IDs starting with numbers like 2512.xxxxx)
	arxivID := extractArxivID(paperLink)
	if arxivID == "" {
		return "" // Not an arXiv paper
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("[POSTER ERROR] Failed to create output directory: %v\n", err)
		return ""
	}

	// Check if arxiv2poster is available
	if _, err := exec.LookPath("arxiv2poster"); err != nil {
		fmt.Printf("[POSTER SKIP] arxiv2poster not found in PATH. Install from: https://github.com/deep-diver/arxiv2poster\n")
		return ""
	}

	// Build arxiv2poster command
	// Using pro model for better quality and aspect ratio control
	args := []string{
		arxivID,
		"--output-dir", outputDir,
		"--orientation", "landscape",
		"--model", "pro", // Use pro for better quality and aspect ratio support
		"--resolution", "4K",
	}

	// Run arxiv2poster
	cmd := exec.Command("arxiv2poster", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[POSTER ERROR] Failed to generate poster for %s: %v\nOutput: %s\n", arxivID, err, string(output))
		return ""
	}

	// arxiv2poster saves to: poster_<arxiv_id>_landscape_english_nopanel.png
	// We need to construct the expected filename
	expectedFilename := fmt.Sprintf("poster_%s_landscape_english_nopanel.png", strings.ReplaceAll(arxivID, "/", "_"))
	posterPath := fmt.Sprintf("%s/%s", outputDir, expectedFilename)

	// Verify the file exists
	if _, err := os.Stat(posterPath); err != nil {
		fmt.Printf("[POSTER ERROR] Expected poster file not found: %s\n", posterPath)
		return ""
	}

	fmt.Printf("[POSTER] Generated poster for %s: %s\n", arxivID, posterPath)
	return posterPath
}

// extractArxivID extracts arXiv ID from HuggingFace paper link
// HuggingFace IDs like "2512.13043" correspond to arXiv IDs
func extractArxivID(link string) string {
	// Pattern: https://huggingface.co/papers/2512.13043
	// The ID after /papers/ is the arXiv ID in most cases
	re := regexp.MustCompile(`huggingface\.co/papers/([\d.]+)`)
	matches := re.FindStringSubmatch(link)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// GeneratePostersForArticles generates posters for all articles that have arXiv IDs
func GeneratePostersForArticles(articles []Article, posterDir string) {
	for i := range articles {
		if articles[i].PosterPath == "" { // Only generate if not already set
			posterPath := GeneratePoster(articles[i].Link, posterDir)
			articles[i].PosterPath = posterPath
		}
	}
}
