package internal

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func ParseDate(date string) (time.Time, error) {
	layout := "2006-01-02 15:04"
	parsedDate, err := time.Parse(layout, date)
	return parsedDate, err
}

func ParseArticle(article string) Article {
	var newArticle Article

	if err := yaml.NewDecoder(strings.NewReader(article)).Decode(&newArticle); err != nil {
		fmt.Println(err)
	}

	return newArticle
}
