# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based automated newsletter service that curates daily ML research papers from Hugging Face and emails them to subscribers. The entire system runs serverlessly on GitHub Actions with no paid infrastructure.

**Tech Stack:** Go 1.21.5, Cobra CLI, Google Gemini 2.5 Flash API, Upstage Solar Pro 2 API, Gmail SMTP

## Build and Run Commands

```bash
# Build the project
go build ./src

# Parse papers from HuggingFace API JSON and generate YAML files
go run src/main.go parse -g <GEMINI_API_KEY> -s <SOLAR_API_KEY> -f <input_json> -o <output_dir>

# Send newsletter (uses YAML files from directory, fills templates, sends email)
go run src/main.go publish <source_dir> -s <date>

# Test without sending email (generates HTML and archives, but skips SMTP)
go run src/main.go publish <source_dir> -s <date> --dry-run

# Archive papers (placeholder - currently unimplemented)
go run src/main.go archive
```

## Architecture

### GitHub Actions Workflows

The system operates through two main workflows:

1. **collect_papers.yml** (runs daily at 20:00 UTC):
   - Downloads papers from HuggingFace API for current date
   - Runs `parse` command to generate YAML files in `current/`
   - Commits changes back to repo

2. **newsletter.yml** (triggered after translation-ko workflow completes):
   - Runs `publish` command on `current/` directory
   - Archives papers to `archive/` directory
   - Generates tag-based markdown indexes
   - Sends HTML email to subscribers
   - Cleans up `current/` directory

3. **translation-ko.yml** and **mov_to_gif.yml**:
   - Korean translation workflow (triggers newsletter)
   - Video to GIF conversion for paper thumbnails

### Code Structure

```
src/
├── main.go                 # Entry point - calls cmd.Execute()
├── cmd/                    # Cobra CLI command definitions
│   ├── root.go            # Base command
│   ├── parse.go           # Parse HuggingFace JSON → YAML files (calls Gemini 2.5 Flash / Solar Pro 2 APIs)
│   ├── publish.go         # Send newsletter (placeholder, actual logic elsewhere)
│   └── archive.go         # Archive management (placeholder)
└── internal/              # Business logic shared across commands
    ├── solar.go           # Solar Pro 2 API client for abstract summarization
    ├── email.go           # Email sending via Gmail SMTP
    ├── newsletter.go      # Data structures (Article, Newsletter)
    ├── configs.go         # Configuration loading
    ├── templ_filler.go    # Template filling logic
    └── utils.go           # Utilities (GetTemplatesInDir, etc.)
```

### Key Data Flow

1. **Parse Phase**: HuggingFace JSON → Go structs → Gemini 2.5 Flash (categorization) + Solar Pro 2 (summarization) → YAML files in `current/`

2. **Publish Phase**: YAML files → Go template (templates/hf_newsletter_template.gohtml) → HTML email → SMTP send

3. **Archive Phase**: YAML files → `archive/<paper-id-first-3-chars>/` + tag index markdown files

### Required Secrets (GitHub Actions)

- `GEMINI_API_KEY` - For paper categorization (20 predefined ML categories) - uses Gemini 2.5 Flash
- `SOLAR_API_KEY` - Upstage Solar Pro 2 API for abstract summarization
- `EMAIL` - Gmail address for sending
- `EMAIL_PASSWORD` - Gmail app-specific password

### Configuration

`configs.yaml` contains:
- Email settings (title, footer, receivers list)
- Header configuration (title, description, image URL)
- Styling (bgcolor, button titles)
- GitHub repository info

### Data Structures

The `Article` struct (internal/newsletter.go) is the core data model:
```go
type Article struct {
    Date            string            // YYYY-MM-DD
    Title           string
    Author          string
    Thumbnail       string            // Media URL
    Link            string            // HuggingFace paper URL
    Tags            []string          // ML categories from Gemini
    Summary         string            // SOLAR-generated summary
    Opinion         string            // Placeholder for future use
    TranslatedPaths map[string]string // Language code → translated file path
}
```

### Paper Categorization

Uses Gemini 2.5 Flash via `SuggestCategories()` in cmd/parse.go - 20 ML categories including:
- Supervised/Unsupervised/Reinforcement Learning
- NLP, Computer Vision, Speech Recognition
- Recommender Systems, Bioinformatics, Robotics
- Fairness/Ethics, Explainable AI, etc.

### Email Templates

Located in `templates/` directory. The main template is `hf_newsletter_template.gohtml` which uses Go's `text/template` with custom functions (`add`, `sub` for index arithmetic). Template filling logic in internal/templ_filler.go pairs articles into tuples for two-column layout.

### Archive Structure

Papers are archived under `archive/<first-3-chars-of-paper-id>/` with tag indexes generated as markdown files linking to archived YAMLs.
