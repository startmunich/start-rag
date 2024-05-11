package main

import (
	"context"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/meilisearch/meilisearch-go"
	"log"
	"notioncrawl/services/api"
	"notioncrawl/services/crawler"
	"notioncrawl/services/crawler/content_crawler/unofficial_content_crawler"
	"notioncrawl/services/crawler/meta_crawler/unofficial_meta_crawler"
	"notioncrawl/services/crawler/workspace_exporter/unofficial_workspace_exporter"
	"notioncrawl/services/notion"
	"notioncrawl/services/state"
	"notioncrawl/services/utils/run_mgr"
	"notioncrawl/services/vector_queue"
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
	tokenv2 := mustEnv("TOKEN_V2")
	spaceId := mustEnv("SPACE_ID")
	startPageId := mustEnv("START_PAGE_ID")
	reRunDelaySec := mustParseInt64(mustEnv("RERUN_DELAY_SEC"))

	meilisearchUrl := mustEnv("MEILISEARCH_URL")
	meilisearchApiToken := mustEnv("MEILISEARCH_API_TOKEN")
	vectorQueueUrl := mustEnv("VECTOR_QUEUE_URL")

	neo4jUrl := mustEnv("NEO4J_URL")
	neo4jUser := mustEnv("NEO4J_USER")
	neo4jPass := mustEnv("NEO4J_PASS")

	influxDbUrl := mustEnv("INFLUXDB_URL")
	influxDbToken := mustEnv("INFLUXDB_TOKEN")
	influxDbOrg := mustEnv("INFLUXDB_ORG")
	influxDbBucket := mustEnv("INFLUXDB_BUCKET")

	port := mustEnv("PORT")

	corsDomains := mustEnv("CORS")

	reRunDelayDuration := time.Second * time.Duration(reRunDelaySec)

	vectorQueue := vector_queue.New(vectorQueueUrl)

	influxDb := influxdb2.NewClient(influxDbUrl, influxDbToken)
	influxWriteAPI := influxDb.WriteAPIBlocking(influxDbOrg, influxDbBucket)

	neo4jOptions := crawler.Neo4jOptions{
		Address:  neo4jUrl,
		Username: neo4jUser,
		Password: neo4jPass,
	}

	meiliClient := meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   meilisearchUrl,
		APIKey: meilisearchApiToken,
	})

	meiliIndex := meiliClient.Index("pages")

	downloadDir, err := os.MkdirTemp("", "notioncrawler_download")
	if err != nil {
		panic("Failed to create temp download folder")
	}
	defer os.RemoveAll(downloadDir)

	notionClient := notion.New(notion.Options{
		NotionSpaceId: spaceId,
		Token:         tokenv2,
		DownloadDir:   downloadDir,
	})

	stateMgr := state.New()
	runMgr := run_mgr.New()

	go api.Run(stateMgr, runMgr, neo4jOptions, meiliIndex, vectorQueue, fmt.Sprintf(":%s", port), corsDomains)

	println("Waiting for Vector Queue ...")
	vectorQueue.WaitForReady()
	println("Vector Queue is ready")

	// USE Exporter to crawl children
	metaCrawler := unofficial_meta_crawler.New(notionClient)
	childrenCrawler := unofficial_content_crawler.New(notionClient)
	workspaceExporter := unofficial_workspace_exporter.New(notionClient)

	for {
		start := time.Now()
		runMgr.Reset()
		log.Printf("Starting Notioncrawler")
		if err := influxWriteAPI.WritePoint(context.Background(), influxdb2.NewPointWithMeasurement("notion_crawler_started").
			SetTime(time.Now())); err != nil {
			log.Println("Failed to write influxdb point")
		}
		stateMgr.UpdateIsRunning(true).UpdateLastRunStartedAt(time.Now().UTC().UnixMilli())
		crawlerInstance := crawler.New(
			stateMgr,
			neo4jOptions,
			vectorQueue,
			meiliIndex,
			startPageId,
			metaCrawler,
			childrenCrawler,
			workspaceExporter,
			&crawler.Options{
				ForceUpdateAll: false,
				ForceUpdateIds: []string{},
			},
		)

		processed := uint64(0)
		cacheMisses := uint64(0)
		wasCanceled := false
		for crawlerInstance.HasNext() {
			if runMgr.ShouldCancel() {
				log.Println(fmt.Sprintf("Run Canceled!"))
				wasCanceled = true
				break
			}

			log.Println(fmt.Sprintf("Queue Size: %d", crawlerInstance.QueueSize()))
			stateMgr.UpdateInQueue(uint64(crawlerInstance.QueueSize())).UpdateProcessed(processed).UpdateCacheMisses(cacheMisses)

			if err := influxWriteAPI.WritePoint(context.Background(), influxdb2.NewPointWithMeasurement("notion_crawler_processing_item_started").
				AddField("processed", processed).
				AddField("cacheMisses", cacheMisses).
				AddField("timeElapsed", time.Since(start).Milliseconds()).
				SetTime(time.Now())); err != nil {
				log.Println("Failed to write influxdb point")
			}

			startItem := time.Now()
			res, err := crawlerInstance.CrawlNext()
			if err != nil {
				log.Println(err.Error())
			} else if res.CacheMiss {
				cacheMisses += 1
			}
			elapsedItem := time.Since(startItem)
			processed += 1

			if err := influxWriteAPI.WritePoint(context.Background(), influxdb2.NewPointWithMeasurement("notion_crawler_processing_item_ended").
				AddField("processed", processed).
				AddField("cacheMisses", cacheMisses).
				AddField("timeElapsed", time.Since(start).Milliseconds()).
				AddField("wasCacheMiss", res.CacheMiss).
				AddField("itemCrawlDuration", elapsedItem.Milliseconds()).
				SetTime(time.Now())); err != nil {
				log.Println("Failed to write influxdb point")
			}
		}
		crawlerInstance.Close()

		elapsed := time.Since(start)
		log.Printf("Notioncrawler took %s", elapsed)
		if err := influxWriteAPI.WritePoint(context.Background(), influxdb2.NewPointWithMeasurement("notion_crawler_ended").
			AddField("processed", processed).
			AddField("cacheMisses", cacheMisses).
			AddField("timeElapsed", elapsed.Milliseconds()).
			AddField("wasCanceled", wasCanceled).
			SetTime(time.Now())); err != nil {
			log.Println("Failed to write influxdb point")
		}
		stateMgr.UpdateIsRunning(false).UpdateLastRunDuration(
			uint64(elapsed.Milliseconds()),
		).UpdateLastRunEndedAt(
			time.Now().UTC().UnixMilli(),
		).UpdateNextRunAt(
			time.Now().UTC().UnixMilli() + reRunDelayDuration.Milliseconds(),
		).UpdateInQueue(
			0,
		).UpdateProcessed(processed).UpdateCacheMisses(cacheMisses)
		time.Sleep(reRunDelayDuration)
	}
}
