/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"context"
	"strings"
    "encoding/json"
    "io/ioutil"
    "os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/deep-diver/hf-daily-paper-newsletter/internal"
)

var Filename string
var OutputDir string
var GeminiAPIKey string
var SolarAPIKey string
var WithPosters bool

type Author struct {
	Id string `json:"_id"`
	Name string `json:"name"`
	Hidden bool `json:"hidden"`
}

type PaperDetail struct {
	Id string `json:"id"`
	Authors []Author `json:"authors"`
	PublishedAt string `json:"publishedAt"`
	Title string `json:"title"`
	Summary string `json:"summary"`
	Upvotes int `json:"upvotes"`
}

type DailyPaper struct {
	Paper PaperDetail `json:"paper"`
    PublishedAt string `json:"publishedAt"`
    Title string `json:"title"`
	Thumbnail string `json:"thumbnail`
    NumComments int `json:"numComments"`
}

type PaperPost struct {
	Date string `yaml:"date"`
	Author string `yaml:"author"`
	Title string `yaml:"title"`
	Thumbnail string `yaml:"thumbnail"`
	Link string `yaml:"link"`
	Summary string `yaml:"summary"`
	Opinion string `yaml:"opinion"`
	Tags []string `yaml:"tags"`
	PosterPath string `yaml:"poster_path,omitempty"`
}

// publishCmd represents the publish command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse HuggingFace daily papers and write yaml files",
	Long: "parse HuggingFace daily papers and write yaml files",

	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		outputdir, _ := cmd.Flags().GetString("outputdir")
		geminiapikey, _ := cmd.Flags().GetString("geminiapikey")
		withPosters, _ := cmd.Flags().GetBool("with-posters")

		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		jsonBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		var daily_papers []DailyPaper
		err = json.Unmarshal(jsonBytes, &daily_papers)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		for _, paper := range daily_papers {
			title := strings.Replace(paper.Title, "\n", " ", -1)
			abstract := strings.Replace(paper.Paper.Summary, "\n", " ", -1)
			summary := internal.SummarizeAbstractWithGemini(geminiapikey, title, abstract)
			tags := SuggestCategories(geminiapikey, title, summary)

			paperLink := "https://huggingface.co/papers/" + paper.Paper.Id
			posterPath := ""

			// Generate poster if experimental flag is enabled
			if withPosters {
				fmt.Printf("[POSTER] Generating poster for %s...\n", paper.Paper.Id)
				posterPath = internal.GeneratePoster(paperLink, "posters")
			}

			paper_post := PaperPost {
				Date: paper.PublishedAt[:10],
				Author: paper.Paper.Authors[0].Name,
				Title: title,
				Thumbnail: paper.Thumbnail,
				Link: paperLink,
				Summary: summary + "...",
				Opinion: "placeholder",
				Tags: tags,
				PosterPath: posterPath,
			}

			data, err := yaml.Marshal(&paper_post)
			if err != nil {
				fmt.Printf("Error marshalling to YAML: %v", err)
				return
			}

			out_filepath := outputdir + "/" + paper.PublishedAt[:10] + " " + title + ".yaml"
			err = os.WriteFile(out_filepath, data, 0644)
			if err != nil {
				fmt.Printf("Error writing YAML file: %v", err)
				return
			}
		}	
	},
}

func SuggestCategories(geminiapikey string, title string, abstract string) []string {
	prompt := `Based on the following information of an academic research paper, suggest EXACTLY 2 specific tags that best describe the paper's focus.

IMPORTANT: Avoid overly generic terms like "Machine Learning", "Deep Learning", "Computer Vision", "Natural Language Processing", "Artificial Intelligence". Instead, use specific sub-fields, techniques, or applications.

Examples of GOOD specific tags:
- Instead of "Computer Vision" → use "Diffusion Models", "Object Detection", "Image Segmentation", "Visual Reasoning", "Video Understanding"
- Instead of "NLP" → use "Transformers", "Large Language Models", "Question Answering", "Machine Translation", "Retrieval Augmented Generation"
- Instead of "Deep Learning" → use "Graph Neural Networks", "Federated Learning", "Prompt Engineering", "Model Compression"
- Instead of generic terms → use "Multimodal Learning", "Reinforcement Learning", "Knowledge Distillation", "Few-shot Learning", "Embodied AI"

Title: "%s"
Abstract: "%s"

Return EXACTLY 2 specific tags as JSON: {"categories": ["tag1", "tag2"]}
	`
	prompt = fmt.Sprintf(prompt, title, abstract)

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiapikey))
	if err != nil {
		fmt.Printf("[TAG ERROR] Client creation failed for '%s': %v\n", title[:50], err)
		return []string{"ML"}
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		fmt.Printf("[TAG ERROR] GenerateContent failed for '%s': %v\n", title[:50], err)
		return []string{"ML"}
	}

	resp_text := fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
	tags := parseResponse(resp_text)
	// Limit to 2 tags
	if len(tags.Categories) > 2 {
		tags.Categories = tags.Categories[:2]
	}
	return tags.Categories
}

type Tags struct {
    Categories []string `json:"categories"`
}

func parseResponse(str string) Tags {
    startIndex := strings.Index(str, "{")
    endIndex := strings.LastIndex(str, "}")

    // Check if both "{" and "}" are found
    if startIndex == -1 || endIndex == -1 {
        fmt.Printf("[TAG ERROR] No JSON found in response: %s\n", str[:200])
        return Tags {Categories: []string{"ML"}}
    }

    // Extract and return the substring
    // Add 1 to endIndex to include the "}" in the result
    trimmedResp := str[startIndex : endIndex+1]

    // Variable to hold the unmarshalled data
    var tags Tags

    // Unmarshal the JSON into the struct
    err := json.Unmarshal([]byte(trimmedResp), &tags)
    if err != nil {
        fmt.Printf("[TAG ERROR] JSON parse failed: %v | Response: %s\n", err, trimmedResp[:200])
        return Tags {Categories: []string{"ML"}}
    }

	return tags
}

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringVarP(&Filename, "filename", "f", "", "filename")
	parseCmd.Flags().StringVarP(&OutputDir, "outputdir", "o", "", "outputdir")
	parseCmd.Flags().StringVarP(&GeminiAPIKey, "geminiapikey", "g", "", "geminiapikey")
	parseCmd.Flags().StringVarP(&SolarAPIKey, "solarapikey", "s", "", "solarapikey")
	parseCmd.Flags().BoolVar(&WithPosters, "with-posters", false, "experimental: generate poster images using arxiv2poster")
}
