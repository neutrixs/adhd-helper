package handlers

import (
	"path/filepath"
	"strings"

	"adhd-helper/models"

	"github.com/gofiber/fiber/v2"
)

// SearchTopics handles recursive full-text search
func SearchTopics(c *fiber.Ctx) error {
	query := c.Query("q")
	scope := c.Query("scope", "")

	if query == "" {
		return c.JSON([]models.SearchResult{})
	}

	// Split query into individual terms
	terms := strings.Fields(strings.TrimSpace(query))
	if len(terms) == 0 {
		return c.JSON([]models.SearchResult{})
	}

	// Determine search directory
	searchDir := ContentDir
	if scope != "" {
		scope = filepath.Clean(scope)
		searchDir = filepath.Join(ContentDir, scope)
	}

	results := models.SearchFiles(searchDir, terms, ContentDir, 20)

	if results == nil {
		results = []models.SearchResult{}
	}

	return c.JSON(results)
}
