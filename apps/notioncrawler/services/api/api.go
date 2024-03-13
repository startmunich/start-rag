package api

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"notioncrawl/services/state"
)

func Run(state *state.Manager, addr string) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Notion Crawler")
	})

	app.Get("/state", func(c *fiber.Ctx) error {
		s, err := json.Marshal(state.GetState())
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		return c.Send(s)
	})

	log.Fatal(app.Listen(addr))
}
