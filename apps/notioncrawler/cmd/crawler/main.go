package main

import (
	"context"
	"fmt"
	"log"
	"notioncrawl/services/crawler"
	"notioncrawl/services/crawler/content_crawler/unofficial_content_crawler"
	"notioncrawl/services/crawler/meta_crawler/unofficial_meta_crawler"
	"notioncrawl/services/crawler/workspace_exporter/unofficial_workspace_exporter"
	"notioncrawl/services/notion"
	"notioncrawl/services/vectordb"
	"os"
	"strconv"
	"time"
)

func mustEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	panic(fmt.Sprintf("Missing required environment variable '%s'", key))
}

func mustParseInt64(num string) uint64 {
	if i, err := strconv.ParseUint(num, 10, 64); err != nil {
		panic(fmt.Sprintf("Cannot parse int64 of '%s'", num))
	} else {
		return i
	}
}

func main() {
	start := time.Now()

	tokenv2 := mustEnv("TOKEN_V2")
	spaceId := mustEnv("SPACE_ID")
	startPageId := mustEnv("START_PAGE_ID")
	qdrantHost := mustEnv("QDRANT_HOST")

	vectorSize := mustEnv("VECTOR_SIZE")
	vectorDistance := mustEnv("VECTOR_DISTANCE")
	vectorCollectionName := mustEnv("VECTOR_COLLECTION_NAME")

	vectorDbOptions := vectordb.QdrantDbOptions{
		Address:        qdrantHost,
		Size:           mustParseInt64(vectorSize),
		Distance:       vectordb.MustParseDistance(vectorDistance),
		CollectionName: vectorCollectionName,
	}

	{
		// Database Setup
		qdClient, err := vectordb.New(vectorDbOptions)
		if err != nil {
			log.Fatal(err)
		}
		qdClient.Setup(context.Background())
		qdClient.Close()
	}

	notionClient := notion.New(notion.Options{
		NotionSpaceId: spaceId,
		Token:         tokenv2,
	})

	// USE Exporter to crawl children
	metaCrawler := unofficial_meta_crawler.New(notionClient)
	childrenCrawler := unofficial_content_crawler.New(notionClient)
	workspaceExporter := unofficial_workspace_exporter.New(notionClient)

	crawlerInstance := crawler.New(
		vectorDbOptions,
		startPageId,
		metaCrawler,
		childrenCrawler,
		workspaceExporter,
		&crawler.Options{
			ForceUpdateAll: false,
			ForceUpdateIds: []string{},
		},
	)
	defer crawlerInstance.Close()

	// Do full export and memgraph import
	if err := crawlerInstance.PerformFullBaseExport(); err != nil {
		return
	}

	elapsed := time.Since(start)
	log.Printf("Notioncrawler took %s", elapsed)
}
