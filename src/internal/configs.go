package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type EmailConfig struct {
	Title       string   `yaml:"title"`
	FooterTitle string   `yaml:"footer_title"`
	Receivers   []string `yaml:"receivers"`

	HeaderTitle           string `yaml:"header_title"`
	HeaderDescription     string `yaml:"header_description"`
	HeaderImageURL        string `yaml:"header_image_url"`
	CommunityLink         string `yaml:"community_link"`
	CommunityLinkBtnTitle string `yaml:"community_link_btn_title"`

	SectionTitle string `yaml:"section_title"`

	ContentLinkBtnTitle string `yaml:"content_link_btn_title"`

	BgColor string `yaml:"bgcolor"`
}

type GeneralConfig struct {
	GitHubRepoURL string `yaml:"github_repo_url"`
	GitBranch     string `yaml:"git_branch"`
}

// Config Type
type Config struct {
	General GeneralConfig `yaml:"general"`
	Email   EmailConfig   `yaml:"email"`
}

func GetConfigs(filename string) Config {
	filename, _ = filepath.Abs(filename)
	yamlFile, _ := ioutil.ReadFile(filename)

	var config Config

	if err := yaml.NewDecoder(strings.NewReader(string(yamlFile))).Decode(&config); err != nil {
		fmt.Println(err)
	}

	return config
}
