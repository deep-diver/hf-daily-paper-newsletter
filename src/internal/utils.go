package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

// ComputeHighlights generates highlights from a list of articles
func ComputeHighlights(articles []Article) Highlights {
	highlights := Highlights{
		PaperCount: len(articles),
		AllTitles:  make([]string, 0, len(articles)),
	}

	// Collect all titles
	for _, article := range articles {
		highlights.AllTitles = append(highlights.AllTitles, article.Title)
	}

	// Count tag frequencies
	tagCounts := make(map[string]int)
	for _, article := range articles {
		for _, tag := range article.Tags {
			tagCounts[tag]++
		}
	}

	// Convert to slice and sort by count (descending)
	tagCountSlice := make([]TagCount, 0, len(tagCounts))
	for tag, count := range tagCounts {
		tagCountSlice = append(tagCountSlice, TagCount{Tag: tag, Count: count})
	}

	sort.Slice(tagCountSlice, func(i, j int) bool {
		return tagCountSlice[i].Count > tagCountSlice[j].Count
	})

	// Take top 5 tags
	if len(tagCountSlice) > 5 {
		highlights.TopTags = tagCountSlice[:5]
	} else {
		highlights.TopTags = tagCountSlice
	}

	return highlights
}

func GetListOfFilesAt(in string, extension string) []string {
	in, _ = filepath.Abs(in)

	filenames := []string{}
	files, _ := ioutil.ReadDir(in)

	for _, f := range files {
		if filepath.Ext(f.Name()) == extension {
			filenames = append(filenames, fmt.Sprintf("%s/%s", in, f.Name()))
		}
	}

	return filenames
}

func ZipArticleTuples(articles []Article, linkTitle, bgColor string) []ArticleTuple {
	articleTuples := []ArticleTuple{}

	if len(articles) == 1 {
		articleTuples = append(articleTuples,
			ArticleTuple{
				Article1:  articles[0],
				Article2:  Article{},
				LinkTitle: linkTitle,
				BgColor:   bgColor,
			})
		return articleTuples
	}

	for i := 0; i < len(articles); i += 2 {
		if i+1 == len(articles) {
			articleTuples = append(articleTuples,
				ArticleTuple{
					Article1:  articles[i],
					Article2:  Article{},
					LinkTitle: linkTitle,
					BgColor:   bgColor,
				})
		} else {
			articleTuples = append(articleTuples,
				ArticleTuple{
					Article1:  articles[i],
					Article2:  articles[i+1],
					LinkTitle: linkTitle,
					BgColor:   bgColor,
				})
		}
	}

	return articleTuples
}
