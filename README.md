[![ci](https://github.com/codingpot/newsletter_awesome_articles/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/codingpot/newsletter_awesome_articles/actions/workflows/ci.yml) [![send_newsletter](https://github.com/codingpot/newsletter_awesome_articles/actions/workflows/newsletter.yml/badge.svg)](https://github.com/codingpot/newsletter_awesome_articles/actions/workflows/newsletter.yml)

# Publish Newsletter Curated by a Group of People

Even though the name says **Group of People**, it can be just you. The aim of this project is to publish and archive newsletters to a target email address.

1. Write a content to be included in newsletter in YAML format under `current` directory.
   - Examples can be found in [test](https://github.com/codingpot/newsletter_awesome_articles/tree/main/test) directory.
2. Create a PR for the write-up and merge it.
3. If merged PRs exceeds the number of **N**, a newsletter contains those contents will be published.
4. Published newsletters will be archived under `archive` directory by assigning the issue number.
5. Also they will be archived by assigned tags (Markdown for each tag will be created, then each markdown contains a list of links to the original YAML file)

![](https://i.ibb.co/2WT5H4P/flow1.png)
![](https://i.ibb.co/Tcr8RFL/flow2.png)

## Prerequisite

You need to create GitHub Secrets named `EMAIL` and `EMAIL_PASSWORD` with appropriate values. Those values will be used in `.github/workflows/newsletter.yml` GitHub Action file to send out emails.

## How to publish?

It basically collects every yaml files under `/current` directory. The name of yaml file should be formatted as `YYYY-MM-DD Title` so the files can be ordered correctly by themselves. The CLI below shows how to publish a newsletter manually.

```shell
# under src/ directory
$ go run main.go publish current
```

If you want to publish a newsletter based on GitHub Action, you need to configure `.github/workflows/newsletter.yml`. By default, it publish a newsletter when there are **four** yaml files in `current` directory. However, you can change the number of yaml files to trigger the publishing behaviour. Just change every line(below) of `newsletter.yml` as you like.

```yml
steps.number_check.outputs.number == HOW_MANY_YAML_FILES
```

## YAML format

```yaml
date: YYYY-MM-DD hh:mm
author: author name
title: title
thumbnail: image URL
link: link URL to the original
summary: preferably up to 2-3 sentences
opinion: preferably up to 5-8 sentences
tags: ["tag1", "tag2"]
```

## Newsletter layout

The referenced newsletter layout is below that I have used for other purposes.

![](https://i.ibb.co/nzSWpSW/Screen-Shot-2022-01-24-at-11-44-48-PM.png)

## Todo

- [x] Parsing YAML
- [x] Filling Template
- [x] Sending Email
- [x] Move Current YAMLs to Archive
- [x] Write CI/CD script (GitHub Action)
