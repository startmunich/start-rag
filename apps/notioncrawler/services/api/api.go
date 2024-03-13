package api

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func Run() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Notion Crawler")
	})

	log.Fatal(app.Listen(":3000"))
}
