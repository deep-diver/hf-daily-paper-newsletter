package internal

import (
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestAllTemplates(t *testing.T) {
	email := Email{
		Title:       "코딩맛집 뉴스레터",
		FooterTitle: "코딩맛집",
		Header: Head{
			Title:                 "코딩맛집에서 발행하는 뉴스레터 입니다",
			Description:           "이것은 뭥미?",
			ImageURL:              "http://...",
			CommunityLink:         "https://discord.gg/HGPnfzDdkG",
			CommunityLinkBtnTitle: "커뮤니티 둘러보고 가입하러 가기",
		},
		FirstSection: Section{
			Title: "이거슨.... 첫 번째 섹션 타이틀",
		},
		ArticleTuples: []ArticleTuple{
			{
				Article1: Article{
					Title:     "My Article",
					Author:    "Chansung",
					Thumbnail: "https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png",
					Link:      "https://github.com/codingpot/newsletter_awesome_articles",
					Tags:      []string{"blog", "oss"},
					Summary:   "Coding Pot Newsletter Platform",
					Opinion:   "Looks amazing!",
				},
				Article2: Article{
					Title:     "My Article",
					Author:    "Chansung",
					Thumbnail: "https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png",
					Link:      "https://github.com/codingpot/newsletter_awesome_articles",
					Tags:      []string{"blog", "oss"},
					Summary:   "Coding Pot Newsletter Platform",
					Opinion:   "Looks amazing!",
				},
			},
		},
	}

	f, err := os.Create("output.html")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	templateFilenames := GetTemplatesInDir("../../templates")
	assert.Equal(t, 5, len(templateFilenames), "wrong length")

	tmpl, _ := template.ParseFiles(templateFilenames...)
	tmpl.ExecuteTemplate(f, "template.gohtml", email)
}
