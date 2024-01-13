package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetListOfFiles(t *testing.T) {
	creatEmptyFile("../../current/empty1.yaml")
	creatEmptyFile("../../current/empty2.yaml")
	creatEmptyFile("../../current/empty3.yaml")
	creatEmptyFile("../../current/empty4.yaml")
	creatEmptyFile("../../current/empty5.yaml")

	filenames := GetListOfFilesAt("../../current", ".yaml")
	assert.Equal(t, 5, len(filenames), "not enough files")

	for _, filename := range filenames {
		os.Remove(filename)
	}
}

func TestGenerateArticleTuples(t *testing.T) {
	articles := []Article{}
	articles = append(articles, Article{})
	articleTuples := ZipArticleTuples(articles, "게시글로 이동", "orange")
	assert.Equal(t, 1, len(articleTuples))

	articles = append(articles, Article{})
	articleTuples = ZipArticleTuples(articles, "게시글로 이동", "orange")
	assert.Equal(t, 1, len(articleTuples))

	articles = append(articles, Article{})
	articleTuples = ZipArticleTuples(articles, "게시글로 이동", "orange")
	assert.Equal(t, 2, len(articleTuples))

	articles = append(articles, Article{})
	articleTuples = ZipArticleTuples(articles, "게시글로 이동", "orange")
	assert.Equal(t, 2, len(articleTuples))

	articles = append(articles, Article{})
	articleTuples = ZipArticleTuples(articles, "게시글로 이동", "orange")
	assert.Equal(t, 3, len(articleTuples))

	articles = append(articles, Article{})
	articleTuples = ZipArticleTuples(articles, "게시글로 이동", "orange")
	assert.Equal(t, 3, len(articleTuples))

	articles = append(articles, Article{})
	articleTuples = ZipArticleTuples(articles, "게시글로 이동", "orange")
	assert.Equal(t, 4, len(articleTuples))
}
