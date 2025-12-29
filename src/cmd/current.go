/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/deep-diver/hf-daily-paper-newsletter/internal"
	"github.com/spf13/cobra"
)

var SubTitle string
var DryRun bool

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// - how to get argument
		subtitle, _ := cmd.Flags().GetString("subtitle")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		config := internal.GetConfigs("../configs.yaml")
		email_config := config.Email
		// article := internal.Article{}
		// 1. get YAML files in current directory
		filenames := internal.GetListOfFilesAt("../current", ".yaml")
		fmt.Println(filenames)
		// 2. open each files & parse them to Article
		articles := []internal.Article{}
		for _, filename := range filenames {
			buf, _ := ioutil.ReadFile(filename)
			articles = append(articles, internal.ParseArticle(string(buf)))
		}
		articleTuples := internal.ZipArticleTuples(articles, email_config.ContentLinkBtnTitle, email_config.BgColor)
		fmt.Println(articleTuples)

		seqNum := 1 // Fixed sequence number since archiving is disabled
		// 3. fill out template
		email := internal.Email{
			Title:       fmt.Sprintf("%s #%d", email_config.Title, seqNum),
			FooterTitle: email_config.FooterTitle,
			Header: internal.Head{
				Title:                 email_config.HeaderTitle,
				Description:           email_config.HeaderDescription,
				ImageURL:              email_config.HeaderImageURL,
				CommunityLink:         email_config.CommunityLink,
				CommunityLinkBtnTitle: email_config.CommunityLinkBtnTitle,
				BgColor:               email_config.BgColor,
			},
			FirstSection: internal.Section{
				Title: email_config.SectionTitle,
			},
			ArticleTuples: articleTuples,
		}
		fmt.Println(email)

		// 4. parse template and generate HTML
		sender_addr := os.Getenv("EMAIL")
		sender_pass := os.Getenv("PASSWORD")
		subject := email_config.Title + "(" + subtitle + ")"
		r := internal.NewRequest(sender_addr, sender_pass, config.Email.Receivers, subject)
		r.ParseOnly("../templates", email)

		// Write HTML to output.html
		f, err := os.Create("output.html")
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		f.WriteString(r.GetBody())
		fmt.Println("HTML written to output.html")

		// 5. send email (skip if dry-run)
		if !dryRun {
			r.Send("../templates", email)
		} else {
			fmt.Println("[DRY-RUN] Skipping email sending")
		}

		// Archiving disabled
	},
}

func init() {
	publishCmd.AddCommand(currentCmd)
	currentCmd.Flags().StringVarP(&SubTitle, "subtitle", "s", "", "subtitle to be concatenated to the main title defined in configs.yaml")
	currentCmd.Flags().BoolVarP(&DryRun, "dry-run", "d", false, "skip email sending (for testing)")
}
