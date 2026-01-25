package main

import (
	"log"
	"os"

	"raider/internal/db"
	"raider/internal/web"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

// getEnv fetches environment variable or returns fallback
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	// Load .env if exists
	_ = godotenv.Load()

	// Configurable values from env
	host := getEnv("WEB_HOST", "0.0.0.0")
	port := getEnv("WEB_PORT", "8080")
	dbPath := getEnv("DB_PATH", "/tmp/raider.db")
	staticPath := getEnv("STATIC_PATH", "/tmp/static/")

	// Initialize database
	db.InitDB(dbPath)

	// Setup template engine
	engine := html.New("./internal/web/templates", ".html")
	engine.Reload(true)

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/static", staticPath)

	web.SetupRoutes(app)

	log.Printf("Server running at http://%s:%s\n", host, port)
	log.Fatal(app.Listen(host + ":" + port))
}
