package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArticle(t *testing.T) {
	article := Article{
		Title:     "My Article",
		Author:    "Chansung",
		Thumbnail: "https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png",
		Link:      "https://github.com/codingpot/newsletter_awesome_articles",
		Tags:      []string{"blog", "oss"},
		Summary:   "Coding Pot Newsletter Platform",
		Opinion:   "Looks amazing!",
	}

	assert.Equal(t, "My Article", article.Title, "Wrong Title")
}
