package notion

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

const (
	getTasksEndpoint    = "https://www.notion.so/api/v3/getTasks"
	enqueueEndpoint     = "https://www.notion.so/api/v3/enqueueTask"
	loadCachedPageChunk = "https://www.notion.so/api/v3/loadCachedPageChunk"
	taskWaitTime        = 3 * time.Second
	taskWaitTimeout     = 3 * time.Hour
	downloadTimeout     = 3 * time.Hour
	requestTimeout      = 10 * time.Second
)

var idRegex = "([0-9a-f]{8})([0-9a-f]{4})([0-9a-f]{4})([0-9a-f]{4})([0-9a-f]{12})"

type Client struct {
	options Options
	client  *http.Client
}

func New(options Options) *Client {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 0,
	}

	nc := &Client{
		options: options,
		client:  client,
	}

	return nc
}

func (c *Client) DownloadToFile(url string, downloadPath string) (*os.File, error) {
	out, err := os.Create(downloadPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	ctx, cancelContext := context.WithTimeout(context.Background(), downloadTimeout)
	defer cancelContext()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "token_v2="+c.options.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) LoadPageBlock(id string) (*PageChunkBlock, error) {
	cachedChunk, err := c.LoadCachedPageChunk(LoadCachedPageChunkOptions{
		BlockId:     id,
		Limit:       1,
		ChunkNumber: 0,
	})
	if err != nil {
		return nil, err
	}

	recordMap, ok := cachedChunk["recordMap"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to find recordMap")
	}

	blocks, ok := recordMap["block"].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to find recordMap.block")
	}

	pageBlockMap, ok := blocks[migrateIdToDashedId(id)].(map[string]interface{})
	if !ok {
		return nil, errors.New("failed to find recordMap.block.value.<blockId>")
	}

	var pageBlock PageChunkBlock
	if err := mapstructure.Decode(pageBlockMap, &pageBlock); err != nil {
		return nil, err
	}

	return &pageBlock, nil
}

func (c *Client) LoadCachedPageChunk(options LoadCachedPageChunkOptions) (map[string]interface{}, error) {
	jsonData := []byte(c.loadCachedPageChunkJson(options))

	var responseJsonNode map[string]interface{}
	if err := c.PostRequest(loadCachedPageChunk, bytes.NewBuffer(jsonData), &responseJsonNode); err != nil {
		return map[string]interface{}{}, err
	}

	_, ok := responseJsonNode["recordMap"].(map[string]interface{})
	if !ok {
		errorName, ok := responseJsonNode["name"].(string)
		if ok && errorName == "UnauthorizedError" {
			log.Println("[ERROR] UnauthorizedError: seems like your token is not valid anymore. Try to log in to Notion again and replace your old token.")
		}
		return map[string]interface{}{}, fmt.Errorf("error name: %v, error message: %v", responseJsonNode["name"], responseJsonNode["message"])
	}

	return responseJsonNode, nil
}

func (c *Client) TriggerExportTask(options ExportOptions) (string, error) {
	jsonData := []byte(c.getTaskJson(options))

	var responseJsonNode map[string]interface{}
	if err := c.PostRequest(enqueueEndpoint, bytes.NewBuffer(jsonData), &responseJsonNode); err != nil {
		return "", err
	}

	taskId, ok := responseJsonNode["taskId"].(string)
	if !ok {
		errorName, ok := responseJsonNode["name"].(string)
		if ok && errorName == "UnauthorizedError" {
			log.Println("[ERROR] UnauthorizedError: seems like your token is not valid anymore. Try to log in to Notion again and replace your old token.")
		}
		return "", fmt.Errorf("error name: %v, error message: %v", responseJsonNode["name"], responseJsonNode["message"])
	}

	return taskId, nil
}

func (c *Client) GetDownloadLink(taskId string) (string, error) {
	postBody := fmt.Sprintf(`{"taskIds": ["%s"]}`, taskId)

	start := time.Now()
	for time.Since(start) < taskWaitTimeout {
		var results Results
		if err := c.PostRequest(getTasksEndpoint, bytes.NewBufferString(postBody), &results); err != nil {
			return "", err
		}

		if len(results.Results) == 0 {
			time.Sleep(taskWaitTime)
			continue
		}

		result := results.Results[0]

		if result.IsFailure() {
			log.Printf("Notion API workspace export returned a 'failure' state. Reason: %s", result.Error)
			return "", nil
		}

		if result.Status != nil {
			if result.Status.ExportUrl != "" {
				log.Printf("Task finished. Pages exported: %d, Status Type: %s",
					result.Status.PagesExported, result.Status.Type)
				return result.Status.ExportUrl, nil
			} else {
				log.Printf("Waiting for task to finish. Pages exported: %d, Status Type: %s",
					result.Status.PagesExported, result.Status.Type)
			}
		}
		time.Sleep(taskWaitTime)
	}

	log.Println("Notion workspace export failed. After waiting 80 minutes, the export status from the Notion API response was still not 'success'")
	return "", nil // or return an error indicating timeout
}
