package grade

import (
  "github.com/gofiber/fiber/v2"
)

func GetGrades(c *fiber.Ctx) error {
  return c.SendString("All grades")
}

func GetGrade(c *fiber.Ctx) error {
  return c.SendString("Single grade")
}

func NewGrade(c *fiber.Ctx) error {
  return c.SendString("New grade")
}

func DeleteGrade(c *fiber.Ctx) error {
  return c.SendString("Delete grade")
}
