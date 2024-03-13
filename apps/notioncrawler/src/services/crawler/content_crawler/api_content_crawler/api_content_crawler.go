package api_content_crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jomei/notionapi"
	"log"
	"notioncrawl/src/services/crawler"
	"regexp"
	"strings"
)

var pageRegex = "https\\:\\/\\/www\\.notion\\.so/[A-Za-z0-9\\-]*([0-9a-f]{32})"

type NotionApiChildrenCrawler struct {
	client *notionapi.Client
}

func NewApiChildCrawler(notionClient *notionapi.Client) crawler.ContentCrawler {
	return &NotionApiChildrenCrawler{
		client: notionClient,
	}
}

func (c *NotionApiChildrenCrawler) CrawlContent(blockId string) (*crawler.ContentCrawlResult, error) {
	var childPages []string
	blockList := []notionapi.BlockID{
		notionapi.BlockID(blockId),
	}

	for i := 0; i < len(blockList); i++ {
		currentBlockId := blockList[i]
		log.Println(fmt.Sprintf("[info] fetching block %d of %d", i, len(blockList)))
		if blocks, err := fetchAllChildBlocks(c.client, currentBlockId); err == nil {
			for _, block := range blocks {
				switch block.(type) {
				case *notionapi.ChildPageBlock:
					pageId := strings.ReplaceAll(block.GetID().String(), "-", "")
					childPages = append(childPages, pageId)
					break

				default:
					linkedPages := extractLinkedPageIds(block)
					childPages = append(childPages, linkedPages...)
					blockList = append(blockList, block.GetID())
				}
			}
		} else {
			return nil, err
		}
	}

	return &crawler.ContentCrawlResult{
		Children: childPages,
		Contents: []*crawler.CrawlContent{},
	}, nil
}

func fetchAllChildBlocks(client *notionapi.Client, blockId notionapi.BlockID) ([]notionapi.Block, error) {
	var blocks []notionapi.Block
	lastBlockId := ""
	hasMore := true

	for hasMore {
		if result, err := client.Block.GetChildren(context.Background(), blockId, &notionapi.Pagination{
			StartCursor: notionapi.Cursor(lastBlockId),
			PageSize:    10,
		}); err == nil {
			hasMore = result.HasMore
			lastBlockId = result.NextCursor
			for _, child := range result.Results {
				blocks = append(blocks, child)
			}
		} else {
			return []notionapi.Block{}, err
		}
	}

	return blocks, nil
}

func extractLinkedPageIds(block notionapi.Block) []string {
	var pageIds []string
	raw, err := json.Marshal(block)
	if err == nil {
		r := regexp.MustCompile(pageRegex)
		matches := r.FindAllStringSubmatch(string(raw), -1)
		for _, match := range matches {
			if match != nil && len(match) > 1 {
				pageIds = append(pageIds, match[1])
			}
		}
	}
	return pageIds
}
