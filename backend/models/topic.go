package models

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type TopicSummary struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsDir       bool   `json:"isDir"`
}

type Breadcrumb struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Path  string `json:"path"`
}

type TopicDetail struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	ContentHTML string         `json:"contentHtml"`
	Children    []TopicSummary `json:"children"`
	Breadcrumbs []Breadcrumb   `json:"breadcrumbs"`
}

type SearchResult struct {
	Path    string `json:"path"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
}

// FrontMatter holds parsed YAML-like front matter from a .md file
type FrontMatter struct {
	Title       string
	Description string
	Body        string
}

// ParseFrontMatter parses a simple --- delimited front matter block
func ParseFrontMatter(content string) FrontMatter {
	fm := FrontMatter{}
	lines := strings.Split(content, "\n")

	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		fm.Body = content
		return fm
	}

	endIdx := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		fm.Body = content
		return fm
	}

	// Parse front matter fields
	for _, line := range lines[1:endIdx] {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "title":
			fm.Title = val
		case "description":
			fm.Description = val
		}
	}

	// Everything after front matter is body
	if endIdx+1 < len(lines) {
		fm.Body = strings.TrimSpace(strings.Join(lines[endIdx+1:], "\n"))
	}

	return fm
}

// ParseMarkdownFile reads and parses a .md file
func ParseMarkdownFile(path string) (FrontMatter, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FrontMatter{}, err
	}
	return ParseFrontMatter(string(data)), nil
}

// ReadTopicSummary reads the summary info from a directory or .md file
func ReadTopicSummary(basePath string, entry os.DirEntry) (TopicSummary, error) {
	name := entry.Name()
	slug := strings.TrimSuffix(name, ".md")

	if entry.IsDir() {
		indexPath := filepath.Join(basePath, name, "_index.md")
		fm, err := ParseMarkdownFile(indexPath)
		if err != nil {
			// No _index.md, use directory name as title
			return TopicSummary{
				Slug:  slug,
				Title: humanize(slug),
				IsDir: true,
			}, nil
		}
		title := fm.Title
		if title == "" {
			title = humanize(slug)
		}
		return TopicSummary{
			Slug:        slug,
			Title:       title,
			Description: fm.Description,
			IsDir:       true,
		}, nil
	}

	// Regular .md file
	if !strings.HasSuffix(name, ".md") || name == "_index.md" {
		return TopicSummary{}, nil
	}

	fm, err := ParseMarkdownFile(filepath.Join(basePath, name))
	if err != nil {
		return TopicSummary{
			Slug:  slug,
			Title: humanize(slug),
		}, nil
	}

	title := fm.Title
	if title == "" {
		title = humanize(slug)
	}

	return TopicSummary{
		Slug:        slug,
		Title:       title,
		Description: fm.Description,
	}, nil
}

// ListChildren lists all child topics (dirs and .md files) in a directory
func ListChildren(dirPath string) ([]TopicSummary, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	var children []TopicSummary
	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden files and _index.md
		if strings.HasPrefix(name, ".") || strings.HasPrefix(name, "_") {
			continue
		}
		// Only include directories and .md files
		if !entry.IsDir() && !strings.HasSuffix(name, ".md") {
			continue
		}

		summary, err := ReadTopicSummary(dirPath, entry)
		if err != nil {
			continue
		}
		if summary.Slug == "" {
			continue
		}
		children = append(children, summary)
	}

	// Sort: directories first, then alphabetically
	sort.Slice(children, func(i, j int) bool {
		if children[i].IsDir != children[j].IsDir {
			return children[i].IsDir
		}
		return children[i].Slug < children[j].Slug
	})

	return children, nil
}

// SearchFiles recursively searches .md files for matching terms
func SearchFiles(dirPath string, terms []string, basePath string, maxResults int) []SearchResult {
	var results []SearchResult
	searchDir(dirPath, terms, basePath, &results, maxResults)
	return results
}

func searchDir(dirPath string, terms []string, basePath string, results *[]SearchResult, maxResults int) {
	if len(*results) >= maxResults {
		return
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if len(*results) >= maxResults {
			return
		}

		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(dirPath, name)

		if entry.IsDir() {
			searchDir(fullPath, terms, basePath, results, maxResults)
			continue
		}

		if !strings.HasSuffix(name, ".md") {
			continue
		}

		searchFile(fullPath, terms, basePath, results)
	}
}

func searchFile(filePath string, terms []string, basePath string, results *[]SearchResult) {
	fm, err := ParseMarkdownFile(filePath)
	if err != nil {
		return
	}

	// Combine title, description, and body for searching
	fullText := strings.ToLower(fm.Title + " " + fm.Description + " " + fm.Body)

	// All terms must match
	for _, term := range terms {
		if !strings.Contains(fullText, strings.ToLower(term)) {
			return
		}
	}

	// Build relative path
	relPath, _ := filepath.Rel(basePath, filePath)
	relPath = strings.TrimSuffix(relPath, ".md")
	relPath = strings.TrimSuffix(relPath, "/_index")

	// Extract snippet
	snippet := extractSnippet(fm.Body, terms[0], 120)

	title := fm.Title
	if title == "" {
		title = humanize(filepath.Base(strings.TrimSuffix(filePath, ".md")))
	}

	*results = append(*results, SearchResult{
		Path:    "/" + relPath,
		Title:   title,
		Snippet: snippet,
	})
}

func extractSnippet(text string, term string, maxLen int) string {
	lower := strings.ToLower(text)
	termLower := strings.ToLower(term)
	idx := strings.Index(lower, termLower)

	if idx == -1 {
		if len(text) > maxLen {
			return text[:maxLen] + "..."
		}
		return text
	}

	// Find a good starting point
	start := idx - 40
	if start < 0 {
		start = 0
	}

	// Move start to the next word boundary
	if start > 0 {
		spaceIdx := strings.IndexByte(text[start:], ' ')
		if spaceIdx != -1 && spaceIdx < 20 {
			start += spaceIdx + 1
		}
	}

	end := start + maxLen
	if end > len(text) {
		end = len(text)
	}

	snippet := text[start:end]
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(text) {
		snippet = snippet + "..."
	}

	return snippet
}

func humanize(slug string) string {
	s := strings.ReplaceAll(slug, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	// Title case
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	var words []string
	for scanner.Scan() {
		w := scanner.Text()
		if len(w) > 0 {
			words = append(words, strings.ToUpper(w[:1])+w[1:])
		}
	}
	return strings.Join(words, " ")
}
