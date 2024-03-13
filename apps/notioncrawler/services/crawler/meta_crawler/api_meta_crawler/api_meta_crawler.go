package api_meta_crawler

import (
	"context"
	"github.com/jomei/notionapi"
	"notioncrawl/src/services/crawler"
)

type ApiMetaCrawler struct {
	client *notionapi.Client
}

func New(notionClient *notionapi.Client) crawler.MetaCrawler {
	return &ApiMetaCrawler{
		client: notionClient,
	}
}

func (c *ApiMetaCrawler) CrawlMeta(id string, parentId string) (*crawler.CrawledPage, error) {
	page, err := c.client.Page.Get(context.Background(), notionapi.PageID(id))
	if err != nil {
		return nil, err
	}

	hash := page.LastEditedTime.UnixMilli()

	return &crawler.CrawledPage{
		PageID:   id,
		ParentID: parentId,
		Url:      page.URL,
		Hash:     hash,
	}, nil
}
