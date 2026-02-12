package handlers

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"adhd-helper/models"

	"github.com/gofiber/fiber/v2"
	"github.com/yuin/goldmark"
)

var ContentDir string

func SetContentDir(dir string) {
	ContentDir = dir
}

// ListTopics returns all top-level topics
func ListTopics(c *fiber.Ctx) error {
	children, err := models.ListChildren(ContentDir)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read topics"})
	}
	return c.JSON(children)
}

// GetTopic returns detail for a specific topic path
func GetTopic(c *fiber.Ctx) error {
	topicPath := c.Params("*")
	if topicPath == "" {
		return ListTopics(c)
	}

	// Clean the path to prevent directory traversal
	topicPath = filepath.Clean(topicPath)
	fullPath := filepath.Join(ContentDir, topicPath)

	// Check if it's a directory
	info, err := os.Stat(fullPath)
	if err == nil && info.IsDir() {
		return getDirectoryTopic(c, fullPath, topicPath)
	}

	// Try as a .md file
	mdPath := fullPath + ".md"
	if _, err := os.Stat(mdPath); err == nil {
		return getFileTopic(c, mdPath, topicPath)
	}

	return c.Status(404).JSON(fiber.Map{"error": "Topic not found"})
}

func getDirectoryTopic(c *fiber.Ctx, dirPath string, relPath string) error {
	// Read _index.md for this directory
	indexPath := filepath.Join(dirPath, "_index.md")
	var contentHTML string
	var title, description string

	fm, err := models.ParseMarkdownFile(indexPath)
	if err == nil {
		title = fm.Title
		description = fm.Description
		if fm.Body != "" {
			var buf bytes.Buffer
			if err := goldmark.Convert([]byte(fm.Body), &buf); err == nil {
				contentHTML = buf.String()
			}
		}
	}

	if title == "" {
		title = humanizePath(relPath)
	}

	children, _ := models.ListChildren(dirPath)
	breadcrumbs := buildBreadcrumbs(relPath)

	return c.JSON(models.TopicDetail{
		Title:       title,
		Description: description,
		ContentHTML: contentHTML,
		Children:    children,
		Breadcrumbs: breadcrumbs,
	})
}

func getFileTopic(c *fiber.Ctx, filePath string, relPath string) error {
	fm, err := models.ParseMarkdownFile(filePath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to read topic"})
	}

	title := fm.Title
	if title == "" {
		title = humanizePath(relPath)
	}

	var contentHTML string
	if fm.Body != "" {
		var buf bytes.Buffer
		if err := goldmark.Convert([]byte(fm.Body), &buf); err == nil {
			contentHTML = buf.String()
		}
	}

	breadcrumbs := buildBreadcrumbs(relPath)

	return c.JSON(models.TopicDetail{
		Title:       title,
		Description: fm.Description,
		ContentHTML: contentHTML,
		Breadcrumbs: breadcrumbs,
	})
}

func buildBreadcrumbs(relPath string) []models.Breadcrumb {
	parts := strings.Split(relPath, "/")
	var crumbs []models.Breadcrumb

	for i, part := range parts {
		path := strings.Join(parts[:i+1], "/")
		crumbs = append(crumbs, models.Breadcrumb{
			Slug:  part,
			Title: humanizePath(part),
			Path:  "/" + path,
		})
	}

	return crumbs
}

func humanizePath(path string) string {
	base := filepath.Base(path)
	s := strings.ReplaceAll(base, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}
