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
		"MATCH (n:CrawledPage)\nDETACH DELETE n",
		map[string]any{}, neo4j.EagerResultTransformer); err != nil {
		return err
	} else {
		return ctx.Status(http.StatusOK).SendString(fmt.Sprintf("Deleted %d nodes", len(result.Records)))
	}
}

func (c *ApiController) GetPagesCount(ctx *fiber.Ctx) error {
	result, err := neo4j.ExecuteQuery(context.Background(), c.neo4j,
		"MATCH (n:CrawledPage)\nRETURN count(n) as count",
		map[string]any{}, neo4j.EagerResultTransformer)

	if err != nil || len(result.Records) < 1 {
		return ctx.JSON(map[string]any{
			"count": 0,
		})
	}

	return ctx.JSON(map[string]any{
		"count": result.Records[0].Get("count"),
	})
}

func (c *ApiController) GetPages(ctx *fiber.Ctx) error {
	result, err := neo4j.ExecuteQuery(context.Background(), c.neo4j,
		"MATCH (n:CrawledPage)-[r]->(m:CrawledPage)\nRETURN n,r,m",
		map[string]any{}, neo4j.EagerResultTransformer)

	if err != nil {
		return err
	}

	var items []interface{}
	for _, record := range result.Records {
		items = append(items, record.AsMap())
	}

	return ctx.JSON(map[string]any{
		"items": items,
	})
}
