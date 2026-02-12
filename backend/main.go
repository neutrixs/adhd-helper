package main

import (
	"log"
	"os"
	"path/filepath"

	"adhd-helper/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Determine content directory
	contentDir := os.Getenv("CONTENT_DIR")
	if contentDir == "" {
		exe, _ := os.Executable()
		contentDir = filepath.Join(filepath.Dir(exe), "content")
		// Fallback to ./content if running with `go run`
		if _, err := os.Stat(contentDir); os.IsNotExist(err) {
			contentDir = "./content"
		}
	}

	absContentDir, err := filepath.Abs(contentDir)
	if err != nil {
		log.Fatal("Failed to resolve content directory:", err)
	}
	handlers.SetContentDir(absContentDir)
	log.Println("Content directory:", absContentDir)

	app := fiber.New(fiber.Config{
		AppName: "ADHD Helper",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET",
		AllowHeaders: "Content-Type",
	}))

	// API routes
	api := app.Group("/api")
	api.Get("/topics", handlers.ListTopics)
	api.Get("/topics/*", handlers.GetTopic)
	api.Get("/search", handlers.SearchTopics)

	// Serve frontend static files in production
	frontendDist := os.Getenv("FRONTEND_DIST")
	if frontendDist == "" {
		frontendDist = "../frontend/build"
	}
	if _, err := os.Stat(frontendDist); err == nil {
		app.Static("/", frontendDist)
		// SPA fallback
		app.Get("/*", func(c *fiber.Ctx) error {
			return c.SendFile(filepath.Join(frontendDist, "index.html"))
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
