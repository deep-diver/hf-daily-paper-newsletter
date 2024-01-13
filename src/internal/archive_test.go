package internal

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func creatEmptyFile(path string) string {
	emptyFile, _ := os.Create(path)
	emptyFile.Close()
	path, _ = filepath.Abs(path)
	return path
}

func writeStringtoFile(to string, content string) {
	f, err := os.OpenFile(to, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	f.WriteString(content)
	f.Close()
}

func lineCounter(fileName string) int {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	num_lines := 0
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		num_lines++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	file.Close()
	return num_lines
}

func TestMoveFiles(t *testing.T) {
	os.RemoveAll("../../current")
	os.RemoveAll("../../archive")
	os.Mkdir("../../current", 0700)
	os.Mkdir("../../archive", 0700)

	path1 := creatEmptyFile("../../current/empty1.yaml")
	path2 := creatEmptyFile("../../current/empty2.yaml")
	path3 := creatEmptyFile("../../current/empty3.yaml")
	path4 := creatEmptyFile("../../current/empty4.yaml")
	path5 := creatEmptyFile("../../current/empty5.yaml")

	filenames := []string{path1, path2, path3, path4, path5}
	parent_dest, _ := filepath.Abs("../../archive")
	base, _ := filepath.Abs("../..")

	dest := MoveFiles(filenames, parent_dest, base+"/")

	base_dest, _ := filepath.Split(dest[0])
	fmt.Println(base_dest)

	files, _ := ioutil.ReadDir(base_dest)
	assert.Equal(t, 5, len(files), "not enough files")

	os.RemoveAll(base_dest)

	_, err := os.Stat(base_dest)
	assert.Equal(t, true, os.IsNotExist(err), "dest directory still exist")
}

func TestSequenceNumbering(t *testing.T) {
	os.RemoveAll("../../current")
	os.RemoveAll("../../archive")
	os.RemoveAll("../../tags")
	os.Mkdir("../../current", 0700)
	os.Mkdir("../../archive", 0700)

	// expect the archive directory exists
	dest, _ := filepath.Abs("../../archive")
	archiveNumber := GetSequenceNumberFromDirs(dest)
	assert.Equal(t, 1, archiveNumber)

	os.Mkdir(dest+"/1", 0700)
	archiveNumber = GetSequenceNumberFromDirs(dest)
	assert.Equal(t, 2, archiveNumber)

	os.Mkdir(dest+"/2", 0700)
	archiveNumber = GetSequenceNumberFromDirs(dest)
	assert.Equal(t, 3, archiveNumber)

	os.Mkdir(dest+"/shit", 0700)
	archiveNumber = GetSequenceNumberFromDirs(dest)
	assert.Equal(t, 3, archiveNumber)

	os.RemoveAll(dest)
	os.Mkdir(dest, 0700)
}

func TestArchiveByTags(t *testing.T) {
	os.RemoveAll("../../current")
	os.RemoveAll("../../archive")
	os.Mkdir("../../current", 0700)
	os.Mkdir("../../archive", 0700)

	mockArticle1 := `
date: 2022-01-01 15:34
author: 박찬성
title: 첫 아티클
thumbnail: https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png
link: https://github.com/codingpot/newsletter_awesome_articles
summary: Coding Pot Newsletter Platform
opinion: Looks amazing!
tags: ["first"]
`

	mockArticle2 := `
date: 2022-01-02 15:34
author: 박찬성
title: 두 번째 아티클
thumbnail: https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png
link: https://github.com/codingpot/newsletter_awesome_articles
summary: Coding Pot Newsletter Platform
opinion: Looks amazing!
tags: ["first", "second"]
`

	mockArticle3 := `
date: 2022-01-03 15:34
author: 박찬성
title: 세 번째 아티클
thumbnail: https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png
link: https://github.com/codingpot/newsletter_awesome_articles
summary: Coding Pot Newsletter Platform
opinion: Looks amazing!
tags: ["first", "second", "third"]
`

	mockArticle4 := `
date: 2022-01-04 15:34
author: 박찬성
title: 네 번째 아티클
thumbnail: https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png
link: https://github.com/codingpot/newsletter_awesome_articles
summary: Coding Pot Newsletter Platform
opinion: Looks amazing!
tags: ["first", "second", "third", "fourth"]
`

	mockArticle5 := `
date: 2022-01-05 15:34
author: 박찬성
title: 다 섯 번째 아티클
thumbnail: https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png
link: https://github.com/codingpot/newsletter_awesome_articles
summary: Coding Pot Newsletter Platform
opinion: Looks amazing!
tags: ["first", "second", "third", "fourth", "fifth"]
`

	writeStringtoFile("../../current/first.yaml", mockArticle1)
	writeStringtoFile("../../current/second.yaml", mockArticle2)
	writeStringtoFile("../../current/third.yaml", mockArticle3)
	writeStringtoFile("../../current/fourth.yaml", mockArticle4)
	writeStringtoFile("../../current/fifth.yaml", mockArticle5)

	articles := []Article{}

	buf, _ := ioutil.ReadFile("../../current/first.yaml")
	articles = append(articles, ParseArticle(string(buf)))

	buf, _ = ioutil.ReadFile("../../current/second.yaml")
	articles = append(articles, ParseArticle(string(buf)))

	buf, _ = ioutil.ReadFile("../../current/third.yaml")
	articles = append(articles, ParseArticle(string(buf)))

	buf, _ = ioutil.ReadFile("../../current/fourth.yaml")
	articles = append(articles, ParseArticle(string(buf)))

	buf, _ = ioutil.ReadFile("../../current/fifth.yaml")
	articles = append(articles, ParseArticle(string(buf)))

	filenames := []string{"../../current/first.yaml",
		"../../current/second.yaml",
		"../../current/third.yaml",
		"../../current/fourth.yaml",
		"../../current/fifth.yaml"}
	parent_dest, _ := filepath.Abs("../../archive")
	base, _ := filepath.Abs("../..")

	archive_destinations := MoveFiles(filenames, parent_dest, base)

	tag_dest, _ := filepath.Abs("../../tags")

	for i, archive_destination := range archive_destinations {
		RecordArticleByTags(articles[i], tag_dest, archive_destination)
	}

	files, _ := ioutil.ReadDir(tag_dest)
	assert.Equal(t, 5, len(files), "not enough files")

	assert.Equal(t, 5, lineCounter(tag_dest+"/first.md"), "not enough number of lines")
	assert.Equal(t, 4, lineCounter(tag_dest+"/second.md"), "not enough number of lines")
	assert.Equal(t, 3, lineCounter(tag_dest+"/third.md"), "not enough number of lines")
	assert.Equal(t, 2, lineCounter(tag_dest+"/fourth.md"), "not enough number of lines")
	assert.Equal(t, 1, lineCounter(tag_dest+"/fifth.md"), "not enough number of lines")

	os.RemoveAll(tag_dest)
	base_dest, _ := filepath.Split(archive_destinations[0])
	os.RemoveAll(base_dest)
}
