/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"strings"
    "encoding/json"
    "io/ioutil"
    "os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Filename string
var OutputDir string

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
	MediaUrl string `json:"mediaUrl`
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
}

// publishCmd represents the publish command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "parse HuggingFace daily papers and write yaml files",
	Long: "parse HuggingFace daily papers and write yaml files",

	Run: func(cmd *cobra.Command, args []string) {
		filename, _ := cmd.Flags().GetString("filename")
		outputdir, _ := cmd.Flags().GetString("outputdir")

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
			summary := strings.Replace(paper.Paper.Summary, "\n", " ", -1)

			paper_post := PaperPost {
				Date: paper.PublishedAt[:10] + " 00:00",
				Author: paper.Paper.Authors[0].Name,
				Title: title,
				Thumbnail: paper.MediaUrl,
				Link: "https://huggingface.co/papers/" + paper.Paper.Id,
				Summary: summary[:500],
				Opinion: "placeholder",
				Tags: []string {"ML"},
			}

			data, err := yaml.Marshal(&paper_post)
			if err != nil {
				fmt.Printf("Error marshalling to YAML: %v", err)
				return
			}
		
			out_filepath := outputdir + "/" + paper.PublishedAt[:10] + " " + title + ".yaml"
			fmt.Println(out_filepath)
			err = os.WriteFile(out_filepath, data, 0644)
			if err != nil {
				fmt.Printf("Error writing YAML file: %v", err)
				return
			}
		}	
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
	parseCmd.Flags().StringVarP(&Filename, "filename", "f", "", "filename")
	parseCmd.Flags().StringVarP(&OutputDir, "outputdir", "o", "", "outputdir")
}
