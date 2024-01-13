package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDateFormat(t *testing.T) {
	testCases := []struct {
		input       string
		expected    time.Time
		expectedErr bool
	}{
		{
			input:       "2022-01-17 15:34",
			expected:    time.Date(2022, time.January, 17, 15, 34, 0, 0, time.UTC),
			expectedErr: false,
		},
		{
			input:       "2022/01/17",
			expectedErr: true,
		},
		{
			input:       "2022-01-17 15:34:02",
			expectedErr: true,
		},
	}

	for _, v := range testCases {
		t.Run(v.input, func(t *testing.T) {
			actual, err := ParseDate(v.input)
			if v.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, v.expected, actual)
		})
	}
}

func TestParsingMarkdown(t *testing.T) {
	mockArticle := `
date: 2022-01-17 15:34
author: 박찬성
title: 첫 아티클
thumbnail: https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png
link: https://github.com/codingpot/newsletter_awesome_articles
summary: Coding Pot Newsletter Platform
opinion: Looks amazing!
tags: ["greeting", "mock", "oh-my"]
`

	expected := Article{
		Date:      "2022-01-17 15:34",
		Author:    "박찬성",
		Title:     "첫 아티클",
		Thumbnail: "https://github.com/codingpot/newsletter_awesome_articles/blob/main/assets/overview.png",
		Link:      "https://github.com/codingpot/newsletter_awesome_articles",
		Summary:   "Coding Pot Newsletter Platform",
		Opinion:   "Looks amazing!",
		Tags:      []string{"greeting", "mock", "oh-my"},
	}

	assert.Equal(t, expected, ParseArticle(mockArticle))
}
