package unofficial_content_crawler

import (
	"notioncrawl/services/crawler"
	"notioncrawl/services/notion"
	"notioncrawl/services/utils"
	"os"
	"path/filepath"
)

type UnofficialContentCrawler struct {
	client *notion.Client
}

func New(client *notion.Client) crawler.ContentCrawler {
	return &UnofficialContentCrawler{
		client: client,
	}
}

func (c *UnofficialContentCrawler) CrawlContent(blockId string) (*crawler.ContentCrawlResult, error) {
	extracted, err := utils.ExportExtracted(c.client, notion.ExportOptions{
		ResourceType:          notion.ResourceTypeBlock,
		BlockId:               blockId,
		ExportType:            notion.ExportTypeMarkdown,
		ExportFiles:           false,
		ExportComments:        false,
		FlattenExportFiletree: false,
	})
	defer os.RemoveAll(extracted)
	if err != nil {
		return nil, err
	}

	var ids []string
	var contents []*crawler.CrawlContent
	if err := filepath.Walk(extracted,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			bytes, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fileContent := string(bytes)

			extractedIds := utils.ExtractLinkedPageIds(fileContent)
			ids = append(ids, extractedIds...)
			contents = append(contents, &crawler.CrawlContent{
				Type:     utils.GetContentType(path),
				FileName: info.Name(),
				Content:  fileContent,
			})
			return nil
		}); err != nil {
		return nil, err
	}

	return &crawler.ContentCrawlResult{
		Children: ids,
		Contents: contents,
	}, nil
}
