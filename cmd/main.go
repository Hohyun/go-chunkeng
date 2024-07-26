package main

import (
	"log"

	"github.com/Hohyun/go-chunkeng/internal/score"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
  "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", hello)

	app.Get("/api/score", score.GetScores)
	app.Get("/api/score/:id", score.GetScore)
	app.Post("/api/score", score.NewScore)
	app.Put("/api/score/:id", score.UpdateScore)
	app.Delete("/api/score/:id", score.DeleteScore)

	app.Get("/api/chaewoom", score.GetChaewooms)
	app.Post("/api/chaewoom", score.NewChaewoom)
	app.Put("/api/chaewoom", score.UpdateChaewoom)
	app.Delete("/api/chaewoom/:id", score.DeleteChaewoom)

	app.Get("/api/class", score.GetClasses)
	app.Get("/api/class/tree", score.GetClassesTree)
	app.Get("/api/class/groups", score.GetClassGroups)
	app.Get("/api/class/teams/:group_id", score.GetClassTeams)
	// app.Get("/api/class-treedata", score.GetClassesTreeData)

	app.Get("/api/member/:class", score.GetMembers)

	app.Get("/api/framework", score.GetFramework)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New())
  app.Use(recover.New())

	setupRoutes(app)

	log.Fatal(app.Listen(":3011"))
}

func hello(c *fiber.Ctx) error {
	// return c.SendString("Hello World 👋! This is great. Isn't it?")
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
		"Content": "Hello World 👋! This is great. Isn't it?",
	})
}
