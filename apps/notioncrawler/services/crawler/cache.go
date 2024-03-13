package crawler

import (
	"context"
	"encoding/json"
	"github.com/fatih/structs"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"strings"
)

func split(s, sep string) []string {
	if len(s) == 0 {
		return []string{}
	}
	return strings.Split(s, sep)
}

func GetCachedPage(driver neo4j.DriverWithContext, id string) *CacheEntry {
	result, err := neo4j.ExecuteQuery(context.Background(), driver,
		"MATCH (n:CrawledPage { page_id: $page_id })\nRETURN n",
		map[string]any{
			"page_id": id,
		}, neo4j.EagerResultTransformer)

	if err != nil || len(result.Records) < 1 {
		return nil
	}

	node, _, err := neo4j.GetRecordValue[neo4j.Node](result.Records[0], "n")
	if err != nil {
		return nil
	}

	pageId, err := neo4j.GetProperty[string](node, "page_id")
	if err != nil {
		return nil
	}

	hash, err := neo4j.GetProperty[int64](node, "hash")
	if err != nil {
		return nil
	}

	childPagesStr, err := neo4j.GetProperty[string](node, "child_pages")
	if err != nil {
		return nil
	}
	childPages := split(childPagesStr, ";")

	crawledPage := CacheEntry{
		PageID:     pageId,
		Hash:       hash,
		ChildPages: childPages,
	}

	return &crawledPage
}

func UpdateCache(driver neo4j.DriverWithContext, crawledPage *CrawledPage, crawlerId string) error {
	crawlPageData := structs.Map(crawledPage)
	crawlPageData["crawler_id"] = crawlerId
	crawlPageData["id"] = crawlPageData["page_id"]
	crawlPageData["child_pages"] = strings.Join(crawledPage.CrawlResult.Children, ";")

	content, err := json.Marshal(crawledPage.CrawlResult.Contents)
	if err != nil {
		return err
	}
	crawlPageData["content"] = string(content)

	_, err = neo4j.ExecuteQuery(context.Background(), driver,
		"MERGE (n:CrawledPage { page_id: $page_id })\n"+
			"ON CREATE SET n.crawlerId=$crawler_id, n.url=$url, n.content=$content, n.child_pages=$child_pages, n.hash=$hash\n"+
			"ON MATCH SET n.crawlerId=$crawler_id, n.url=$url, n.content=$content, n.child_pages=$child_pages, n.hash=$hash\n",
		crawlPageData, neo4j.EagerResultTransformer)
	if err != nil {
		return err
	}

	if crawledPage.ParentID != "" {
		_, err := neo4j.ExecuteQuery(context.Background(), driver,
			"MATCH (p:CrawledPage { page_id: $parent_id })\nMATCH (c:CrawledPage { page_id: $page_id })\n"+
				"CALL {\n  WITH p,c\n  MATCH (p)-[r:LINKS_TO]->(c)\n  DELETE r\n}\n"+
				"CREATE (p)-[r:LINKS_TO]->(c)",
			map[string]any{
				"parent_id": crawledPage.ParentID,
				"page_id":   crawledPage.PageID,
			}, neo4j.EagerResultTransformer)
		if err != nil {
			return err
		}
	}
	return nil
}
