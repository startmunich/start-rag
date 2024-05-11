package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/meilisearch/meilisearch-go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
)

func (c *ApiController) PurgeDb(ctx *fiber.Ctx) error {
	c.runMgr.CancelCurrentRun()
	if err := c.vectorQueue.PurgeQueue(); err != nil {
		println("Warn: cannot purge vector queue")
		return err
	}
	if err := c.vectorQueue.PurgeVectorDb(); err != nil {
		println("Warn: cannot purge vector database")
		return err
	}
	if _, err := c.meiliIndex.DeleteAllDocuments(); err != nil {
		println("Warn: cannot purge meili index")
		return err
	}
	if result, err := neo4j.ExecuteQuery(context.Background(), c.neo4j,
		"MATCH (n:CrawledPage)\nDETACH DELETE n",
		map[string]any{}, neo4j.EagerResultTransformer); err != nil {
		return err
	} else {
		return ctx.Status(http.StatusOK).SendString(fmt.Sprintf("Deleted %d nodes", len(result.Records)))
	}
}

func (c *ApiController) GetPagesCount(ctx *fiber.Ctx) error {
	ctx.Response().Header.Add("Cache-Time", "60")
	result, err := neo4j.ExecuteQuery(context.Background(), c.neo4j,
		"MATCH (n:CrawledPage)\nRETURN count(n) as count",
		map[string]any{}, neo4j.EagerResultTransformer)

	if err != nil || len(result.Records) < 1 {
		return ctx.JSON(map[string]any{
			"count": 0,
		})
	}

	count, exists := result.Records[0].Get("count")
	if !exists {
		return ctx.JSON(map[string]any{
			"count": 0,
		})
	}
	return ctx.JSON(map[string]any{
		"count": count,
	})
}

func (c *ApiController) GetPages(ctx *fiber.Ctx) error {
	ctx.Response().Header.Add("Cache-Time", "600")
	result, err := neo4j.ExecuteQuery(context.Background(), c.neo4j,
		"MATCH (n:CrawledPage)\nRETURN n{.page_id,.url,.child_pages}",
		map[string]any{}, neo4j.EagerResultTransformer)

	if err != nil {
		return err
	}

	var items []interface{}
	for _, record := range result.Records {
		items = append(items, record.AsMap()["n"])
	}

	return ctx.JSON(map[string]any{
		"pages": items,
	})
}

func (c *ApiController) Search(ctx *fiber.Ctx) error {
	searchQuery := ctx.Query("q")
	searchRes, err := c.meiliIndex.Search(searchQuery, &meilisearch.SearchRequest{
		AttributesToHighlight: []string{"*"},
	})
	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.JSON(searchRes)
}
