package main

import (
	"log"

	"github.com/Hohyun/go-chunkeng/internal/score"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", hello)

	app.Get("/api/score", score.GetScores)
	app.Get("/api/score/:id", score.GetScore)
	app.Post("/api/score", score.NewScore)
	app.Delete("/api/score", score.DeleteScore)

	app.Get("/api/class", score.GetClasses)
	app.Get("/api/class/tree", score.GetClassesTree)
	// app.Get("/api/class-treedata", score.GetClassesTreeData)

	app.Get("/api/member/:class", score.GetMembers)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	setupRoutes(app)

	log.Fatal(app.Listen(":3101"))
}

func hello(c *fiber.Ctx) error {
	return c.SendString("Hello World 👋! This is great. Isn't it?")
}
