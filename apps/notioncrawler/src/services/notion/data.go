package notion

import (
	"fmt"
	"regexp"
	"strings"
)

func (c *Client) loadCachedPageChunkJson(o LoadCachedPageChunkOptions) string {
	return fmt.Sprintf(`{
		"page": {
			"id": "%s"
		},
		"limit": %d,
		"cursor": {
			"stack": []
		},
		"chunkNumber": %d,
		"verticalColumns": false
	}`, migrateIdToDashedId(o.BlockId), o.Limit, o.ChunkNumber)
}

func (c *Client) getTaskJson(o ExportOptions) string {
	if o.ResourceType == ResourceTypeSpace {
		return c.exportSpaceTask(o)
	} else if o.ResourceType == ResourceTypeBlock {
		return c.exportBlockTask(o)
	}
	return ""
}

func (c *Client) exportSpaceTask(o ExportOptions) string {
	return fmt.Sprintf(`{
        "task": {
            "eventName": "exportSpace",
            "request": {
                "spaceId": "%s",
                "shouldExportComments": %t,
                "exportOptions": {
                    "exportType": "%s",
                    "flattenExportFiletree": %t,
                    "timeZone": "Europe/Berlin",
                    "locale": "en"
                }
            }
        }
    }`, c.options.NotionSpaceId, o.ExportComments, o.ExportType, o.FlattenExportFiletree)
}

func (c *Client) exportBlockTask(o ExportOptions) string {
	exportFilesJson := ""
	if !o.ExportFiles {
		exportFilesJson = "\n\t\t\t\t\t\"includeContents\": \"no_files\","
	}

	return fmt.Sprintf(`{
        "task": {
            "eventName": "exportBlock",
            "request": {
                "block": {
					"id": "%s",
					"spaceId": "%s"
				},
                "shouldExportComments": %t,
				"recursive": false,
                "exportOptions": {
                    "collectionViewExportType": "all",
                    "exportType": "%s",%s
                    "locale": "en",
                    "timeZone": "Europe/Berlin"
                }
            }
        }
    }`, migrateIdToDashedId(o.BlockId), c.options.NotionSpaceId, o.ExportComments, o.ExportType, exportFilesJson)
}

func migrateIdToDashedId(id string) string {
	r := regexp.MustCompile(idRegex)
	match := r.FindStringSubmatch(id)
	if match != nil {
		return strings.Join(match[1:], "-")
	}
	return id
}
