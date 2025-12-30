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

// DailyInsight represents the AI-generated daily insight
type DailyInsight struct {
	Insight string `json:"insight"`
}

// GenerateDailyInsight creates a cohesive summary of all papers using Gemini
func GenerateDailyInsight(geminiapikey string, articles []Article) string {
	if len(articles) == 0 {
		return "No papers available for today's newsletter."
	}

	// Build a concise prompt with paper info
	var paperInfo strings.Builder
	paperInfo.WriteString("Here are today's ML research papers from Hugging Face:\n\n")

	for i, article := range articles {
		paperInfo.WriteString(fmt.Sprintf("%d. Title: %s\n", i+1, article.Title))
		paperInfo.WriteString(fmt.Sprintf("   Tags: %s\n", strings.Join(article.Tags, ", ")))
		if len(article.Summary) > 0 {
			// Use first part of summary to stay within token limits
			summaryLen := len(article.Summary)
			if summaryLen > 150 {
				summaryLen = 150
			}
			paperInfo.WriteString(fmt.Sprintf("   Summary: %s...\n", article.Summary[:summaryLen]))
		}
		paperInfo.WriteString("\n")
	}

	prompt := `Analyze the following collection of ML research papers and provide a single cohesive paragraph (3-4 sentences) that:
1. Identifies the main themes and trends across today's papers
2. Highlights any interesting connections or patterns between different papers
3. Gives readers a quick sense of what's exciting in ML research today

Write in an engaging but professional tone. Your answer should be formatted in a valid JSON as {"insight": string}

%s
`

	prompt = fmt.Sprintf(prompt, paperInfo.String())

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiapikey))
	if err != nil {
		fmt.Printf("[INSIGHT ERROR] Gemini client error: %v\n", err)
		return fallbackDailyInsight(articles)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		fmt.Printf("[INSIGHT ERROR] Gemini generation error: %v\n", err)
		return fallbackDailyInsight(articles)
	}

	respText := fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
	return parseDailyInsightResponse(respText, articles)
}

func parseDailyInsightResponse(str string, articles []Article) string {
	startIndex := strings.Index(str, "{")
	endIndex := strings.LastIndex(str, "}")

	if startIndex == -1 || endIndex == -1 {
		return fallbackDailyInsight(articles)
	}

	trimmedResp := str[startIndex : endIndex+1]

	var insight DailyInsight
	err := json.Unmarshal([]byte(trimmedResp), &insight)
	if err != nil {
		return fallbackDailyInsight(articles)
	}

	if len(insight.Insight) == 0 {
		return fallbackDailyInsight(articles)
	}

	return insight.Insight
}

func fallbackDailyInsight(articles []Article) string {
	if len(articles) == 0 {
		return "No papers available for today's newsletter."
	}

	if len(articles) == 1 {
		return fmt.Sprintf("Today's newsletter features %d paper exploring %s.", len(articles), articles[0].Tags[0])
	}

	topTag := articles[0].Tags[0]
	return fmt.Sprintf("Today's newsletter features %d papers covering topics like %s and more.", len(articles), topTag)
}
