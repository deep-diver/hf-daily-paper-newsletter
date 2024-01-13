package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

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
