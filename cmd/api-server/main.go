package main

import (
  "log"

  "github.com/gofiber/fiber/v2"
  "github.com/Hohyun/go-chunkeng/grade"
)

func setupRoutes(app *fiber.App) {
  app.Get("/", hello)

  app.Get("/api/grade", grade.GetGrades)
  app.Get("/api/grade/:id", grade.GetGrade)
  app.Post("/api/grade", grade.NewGrade)
  app.Delete("/api/grade", grade.DeleteGrade)
}

func main() {
  app := fiber.New()

  setupRoutes(app)

  log.Fatal(app.Listen(":3000"))
}

func hello(c *fiber.Ctx) error {
  return c.SendString("Hello World 👋! This is great. Isn't it?")
}
