package unofficial_workspace_exporter

import (
	"errors"
	"fmt"
	"log"
	"notioncrawl/src/services/crawler"
	"notioncrawl/src/services/notion"
	"notioncrawl/src/services/utils"
	"os"
	"path/filepath"
)

type UnofficialWorkspaceExporter struct {
	client        *notion.Client
	extractedPath string
	files         []string
	pointer       int
}

func New(client *notion.Client) crawler.WorkspaceExporter {
	return &UnofficialWorkspaceExporter{
		client:  client,
		files:   []string{},
		pointer: 0,
	}
}

func (c *UnofficialWorkspaceExporter) Export() error {
	extracted, err := utils.ExportExtracted(c.client, notion.ExportOptions{
		ResourceType:          notion.ResourceTypeSpace,
		ExportType:            notion.ExportTypeMarkdown,
		ExportFiles:           false,
		ExportComments:        false,
		FlattenExportFiletree: true,
	})
	if err != nil {
		return err
	}

	c.extractedPath = extracted

	if err := filepath.Walk(extracted,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			c.files = append(c.files, path)
			return nil
		}); err != nil {
		return err
	}

	c.pointer = 0
	return nil
}

func (c *UnofficialWorkspaceExporter) getNextPagePath() (string, error) {
	if c.pointer < len(c.files) {
		defer func() {
			c.pointer++
		}()
		return c.files[c.pointer], nil
	}
	return "", errors.New("no items left")
}

func (c *UnofficialWorkspaceExporter) GetNextPage() *crawler.CrawledPage {
	log.Printf("[info] Get Page: %d of %d", c.pointer, len(c.files))
	path, err := c.getNextPagePath()
	if err != nil {
		return nil
	}

	fileName := filepath.Base(path)

	pathPageIds := utils.ExtractPageIdsFromPath(fileName)
	if len(pathPageIds) < 1 {
		return nil
	}
	pageId := pathPageIds[len(pathPageIds)-1]

	var ids []string
	var contents []*crawler.CrawlContent

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	fileContent := string(bytes)

	extractedIds := utils.ExtractLinkedPageIds(fileContent)
	ids = append(ids, extractedIds...)
	contents = append(contents, &crawler.CrawlContent{
		Type:     utils.GetContentType(path),
		FileName: fileName,
		Content:  fileContent,
	})

	hash := int64(0)

	return &crawler.CrawledPage{
		PageID:   pageId,
		ParentID: "",
		Url:      fmt.Sprintf("https://notion.so/%s", pageId),
		Hash:     hash,
		CrawlResult: &crawler.ContentCrawlResult{
			Children: extractedIds,
			Contents: []*crawler.CrawlContent{},
		},
	}
}

func (c *UnofficialWorkspaceExporter) HasNextPage() bool {
	return c.pointer < len(c.files)
}

func (c *UnofficialWorkspaceExporter) Close() error {
	return os.RemoveAll(c.extractedPath)
}
