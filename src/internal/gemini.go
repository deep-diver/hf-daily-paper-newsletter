package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// SummarizeAbstractWithGemini summarizes an abstract using Gemini API
func SummarizeAbstractWithGemini(geminiapikey string, title string, abstract string) string {
	prompt := `Below is an abstract from an academic paper. Summarize the abstract in two sentences. If possible ELI5, and don't say "We introduce..." or "This paper...". Your answer should be formatted in a valid JSON as {"summary": string}

Title: "%s"
Abstract: "%s"
`
	prompt = fmt.Sprintf(prompt, title, abstract)

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiapikey))
	if err != nil {
		fmt.Printf("Gemini client error: %v\n", err)
		return fallbackSummary(abstract)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		fmt.Printf("Gemini generation error: %v\n", err)
		return fallbackSummary(abstract)
	}

	resp_text := fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
	return parseSummaryResponse(resp_text, abstract)
}

func parseSummaryResponse(str string, abstract string) string {
	startIndex := strings.Index(str, "{")
	endIndex := strings.LastIndex(str, "}")

	if startIndex == -1 || endIndex == -1 {
		return fallbackSummary(abstract)
	}

	trimmedResp := str[startIndex : endIndex+1]

	var summary Summary
	err := json.Unmarshal([]byte(trimmedResp), &summary)
	if err != nil {
		return fallbackSummary(abstract)
	}

	if len(summary.Text) >= 500 {
		return summary.Text[:500]
	}
	return summary.Text
}

func fallbackSummary(abstract string) string {
	if len(abstract) >= 500 {
		return abstract[:500]
	}
	return abstract
}
