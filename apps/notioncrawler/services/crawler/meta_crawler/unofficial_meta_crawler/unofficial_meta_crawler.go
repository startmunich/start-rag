package unofficial_meta_crawler

import (
	"fmt"
	"notioncrawl/services/crawler"
	"notioncrawl/services/notion"
)

type UnofficialMetaCrawler struct {
	client *notion.Client
}

func New(client *notion.Client) crawler.MetaCrawler {
	return &UnofficialMetaCrawler{
		client: client,
	}
}

func (c *UnofficialMetaCrawler) CrawlMeta(id string, parentId string) (*crawler.CrawledPage, error) {
	page, err := c.client.LoadPageBlock(id)
	if err != nil {
		return nil, err
	}

	hash := page.Value.LastEditedTime

	return &crawler.CrawledPage{
		PageID:   id,
		ParentID: parentId,
		Url:      fmt.Sprintf("https://notion.so/%s", id),
		Hash:     hash,
	}, nil
}
