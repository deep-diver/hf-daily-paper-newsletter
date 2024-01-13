package internal

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func GetTemplatesInDir(parent_path string) []string {
	parent_path, _ = filepath.Abs(parent_path)
	filenames := []string{}

	files, err := ioutil.ReadDir(parent_path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, parent_path+"/"+file.Name())
		}
	}

	return filenames
}

func TT() {
	tmpl, _ := template.ParseFiles("/workspaces/newsletter_awesome_articles/templates/template.gohtml")
	tmpl.Execute(os.Stdout, nil)
	// tmpl.Execute("header", mockData)
}
