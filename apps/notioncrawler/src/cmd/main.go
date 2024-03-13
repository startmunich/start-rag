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

const (
	token       = "secret_c6vmUvoXD96zIBFolSK0f4wNgu6j1n1MFxZZMDQitGf"
	tokenv2     = ""
	spaceId     = "3abc121b-7cb1-477b-942d-5404c70daf67"
	startPageId = "297883f9bc6c49f4bcf9d42dbe8fe969"
)

func defaultEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

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

	//notionClient := notionapi.NewClient(token)

	/*
		// EXAMPLE EXPORT ONE PAGE
		export, err := exporterInstance.ExportExtracted(exporter.ExportOptions{
			ResourceType:          exporter.ResourceTypeBlock,
			BlockId:               "fc6b85e65b29408bb5297112cf1c76ec",
			ExportType:            exporter.ExportTypeMarkdown,
			ExportFiles:           false,
			ExportComments:        true,
			FlattenExportFiletree: false,
		})
		if err != nil {
			return
		}

		log.Println(export)
	*/

	// USE Notion Block API to crawl Children
	//childrenCrawler := crawler.NewApiChildCrawler(notionClient)

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

	// Loops over pages following links and updates in memgraph if something changed
	/*for crawlerInstance.HasNext() {
		log.Println(fmt.Sprintf("Queue Size: %d", crawlerInstance.QueueSize()))
		err := crawlerInstance.CrawlNext()
		if err != nil {
			log.Println(err.Error())
		}
	}
	crawlerInstance.Print()*/

	// Do full export and memgraph import
	if err := crawlerInstance.PerformFullBaseExport(); err != nil {
		return
	}

	elapsed := time.Since(start)
	log.Printf("Notioncrawler took %s", elapsed)
}
