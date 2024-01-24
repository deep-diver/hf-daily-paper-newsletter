package internal

import (
	"context"
	"fmt"
	"strings"
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
)

func SummarizeAbstract(solarapikey string, title string, abstract string) string {
	// solarapikey, _ := cmd.Flags().GetString("solarapikey")

	userMessage := `below is an abstract from an academic paper. summarize the abstract in two sentences. If possible ELI5, and don't say "We introduce..." or "This paper...". Your answer should be formatted in a valid JSON as {"summary": string}
--------------
Abstract: "%s"
`
	userMessage = fmt.Sprintf(userMessage, title, abstract)

	config := openai.DefaultConfig(solarapikey)		
	config.BaseURL = "https://api.upstage.ai/v1/solar"
	client := openai.NewClientWithConfig(config)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "upstage/solar-1-mini-chat",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userMessage,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		if len(abstract) >= 500 {
			return abstract[:500]	
		} else {
			return abstract
		}
	}

	summary := parseResponse(resp.Choices[0].Message.Content, abstract)
	return summary
}

type Summary struct {
    Text string `json:"summary"`
}

func parseResponse(str string, abstract string) string {
    startIndex := strings.Index(str, "{")
    endIndex := strings.LastIndex(str, "}")

    // Check if both "{" and "}" are found
    if startIndex == -1 || endIndex == -1 {
		if len(abstract) >= 500 {
			return abstract[:500]	
		} else {
			return abstract
		}
    }

    // Extract and return the substring
    // Add 1 to endIndex to include the "}" in the result
    trimmedResp := str[startIndex : endIndex+1]
	
    // Variable to hold the unmarshalled data
    var summary Summary

    // Unmarshal the JSON into the struct
    err := json.Unmarshal([]byte(trimmedResp), &summary)
    if err != nil {
		if len(abstract) >= 500 {
			return abstract[:500]	
		} else {
			return abstract
		}
    }

	if len(summary.Text) >= 500 {
		return summary.Text[:500]
	} else {
		return summary.Text
	}
}
