package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"notioncrawl/src/services/crawler"
	"notioncrawl/src/services/crawler/content_crawler/unofficial_content_crawler"
	"notioncrawl/src/services/crawler/meta_crawler/unofficial_meta_crawler"
	"notioncrawl/src/services/crawler/workspace_exporter/unofficial_workspace_exporter"
	"notioncrawl/src/services/notion"
	"os"
	"time"
)

func mustEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	panic(fmt.Sprintf("Missing required environment variable '%s'", key))
}

func main() {
	start := time.Now()

	tokenv2 := mustEnv("TOKEN_V2")
	spaceId := mustEnv("SPACE_ID")
	startPageId := mustEnv("START_PAGE_ID")
	qdrantHost := mustEnv("QDRANT_HOST")

	dbUri := "bolt://localhost:7687"
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth("username", "password", ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(context.Background())

	notionClient := notion.New(notion.Options{
		NotionSpaceId: spaceId,
		Token:         tokenv2,
	})

	// USE Exporter to crawl children
	metaCrawler := unofficial_meta_crawler.New(notionClient)
	childrenCrawler := unofficial_content_crawler.New(notionClient)
	workspaceExporter := unofficial_workspace_exporter.New(notionClient)

	crawlerInstance := crawler.New(
		driver,
		startPageId,
		metaCrawler,
		childrenCrawler,
		workspaceExporter,
		&crawler.Options{
			ForceUpdateAll: false,
			ForceUpdateIds: []string{},
		},
	)

	// Do full export and memgraph import
	if err := crawlerInstance.PerformFullBaseExport(); err != nil {
		return
	}

	elapsed := time.Since(start)
	log.Printf("Notioncrawler took %s", elapsed)
}
