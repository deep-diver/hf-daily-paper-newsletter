package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func createDirIfNotExist(dest string) bool {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		os.Mkdir(dest, 0700)
		return true
	}

	return false
}

func MoveFiles(filenames []string, to string, base string) []string {
	result := []string{}

	createDirIfNotExist(to)
	seq_num := GetSequenceNumberFromDirs(to)
	target_to := to + "/" + strconv.Itoa(seq_num)
	createDirIfNotExist(target_to)

	for _, filename := range filenames {
		original_filename := filename

		_, filename := filepath.Split(filename)
		target_filename := target_to + "/" + url.QueryEscape(filename)

		err := os.Rename(original_filename, target_filename)
		if err != nil {
			log.Fatal(err)
		} else {
			tokens := strings.Split(target_filename, "/")
			rel_path := strings.Join(tokens[len(tokens)-3:], "/")
			result = append(result, base+rel_path)
		}
	}

	return result
}

func GetSequenceNumberFromDirs(in string) int {
	createDirIfNotExist(in)

	files, err := ioutil.ReadDir(in)
	if err != nil {
		log.Fatal(err)
	}

	count := 1
	for _, f := range files {
		if f.IsDir() {
			if _, err := strconv.Atoi(f.Name()); err == nil {
				count++
			}
		}
	}

	return count
}

func RecordArticleByTags(article Article, to string, where string) {
	createDirIfNotExist(to)
	for _, tag := range article.Tags {
		// If the file doesn't exist, create it, or append to the file
		filename := fmt.Sprintf("%s/%s.md", to, tag)
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		link_string := fmt.Sprintf("- [%s](%s) / %s\n", article.Title, where, article.Date)
		f.WriteString(link_string)
		f.Close()
	}
}
