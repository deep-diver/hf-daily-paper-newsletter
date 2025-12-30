package internal

type Article struct {
	Date      string   `yaml:"date"`
	Title     string   `yaml:"title"`
	Author    string   `yaml:"author"`
	Thumbnail string   `yaml:"thumbnail"`
	Link      string   `yaml:"link"`
	Tags      []string `yaml:"tags"`
	Summary   string   `yaml:"summary"`
	Opinion   string   `yaml:"opinion"`
	TranslatedPaths map[string]string `yaml:"translated_paths"`
	PosterPath string `yaml:"poster_path,omitempty"` // Path to generated poster image (experimental)
}

type Newsletter struct {
	Articles []Article `yaml:"articles"`
}

// TagCount represents a tag and its frequency
type TagCount struct {
	Tag   string
	Count int
}

// Highlights contains daily summary information
type Highlights struct {
	PaperCount    int        // Total number of papers
	TopTags       []TagCount // Top 5 most common tags with counts
	AllTitles     []string   // All paper titles for quick preview
	DailyInsight  string     // AI-generated insight about the day's papers
}
