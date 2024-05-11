package api

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/meilisearch/meilisearch-go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	"notioncrawl/services/api/controller"
	"notioncrawl/services/crawler"
	"notioncrawl/services/state"
	"notioncrawl/services/vector_queue"
	"time"
)

func Run(state *state.Manager, neo4jOptions crawler.Neo4jOptions, meiliIndex *meilisearch.Index, vectorQueue *vector_queue.VectorQueue, addr string, corsDomains string) {
	neo4jdriver, err := neo4j.NewDriverWithContext(neo4jOptions.Address, neo4j.BasicAuth(neo4jOptions.Username, neo4jOptions.Password, ""))
	if err != nil {
		panic(err)
	}

	ctrl := controller.New(neo4jdriver, meiliIndex, vectorQueue)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: corsDomains,
	}))

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("noCache") == "true"
		},
		Expiration:   3 * time.Second,
		CacheControl: true,
	}))

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

	app.Get("/search", ctrl.Search)

	app.Get("/pages", ctrl.GetPages)
	app.Get("/pages/count", ctrl.GetPagesCount)

	app.Post("/db/purge", ctrl.PurgeDb)

	log.Fatal(app.Listen(addr))
}
