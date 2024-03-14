package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
)

func (c *ApiController) PurgeDb(ctx *fiber.Ctx) error {
	if result, err := neo4j.ExecuteQuery(context.Background(), c.neo4j,
		"MATCH (n:CrawledPage)\nDELETE n",
		map[string]any{}, neo4j.EagerResultTransformer); err != nil {
		return err
	} else {
		return ctx.Status(http.StatusOK).SendString(fmt.Sprintf("Deleted %d nodes", len(result.Records)))
	}
}
