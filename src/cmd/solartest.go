/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	openai "github.com/sashabaranov/go-openai"
)

// Define the request body structure
type requestBody struct {
	Model    string `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool `json:"stream"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var SolarAPIKey string

// publishCmd represents the publish command
var solarCmd = &cobra.Command{
	Use:   "solar",
	Short: "A brief description of your command",
	Long: `A `,
	Run: func(cmd *cobra.Command, args []string) {
		solarapikey, _ := cmd.Flags().GetString("solarapikey")

		client := openai.NewClient(solarapikey)
		client.BaseURL = "https://api.upstage.ai/v1/solar"
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: "upstage/solar-1-mini-chat",
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: "Hello!",
					},
				},
			},
		)
	
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			return
		}
	
		fmt.Println(resp.Choices[0].Message.Content)
	},
}

func init() {
	rootCmd.AddCommand(solarCmd)
	parseCmd.Flags().StringVarP(&SolarAPIKey, "solarapikey", "s", "", "solarapikey")
}
