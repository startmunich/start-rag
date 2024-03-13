package utils

import (
	"notioncrawl/services/crawler"
	"path/filepath"
	"regexp"
)

const pageIdLinkRegex = "https\\:\\/\\/www\\.notion\\.so/[A-Za-z0-9\\-]*([0-9a-f]{32})"
const pageIdPathRegex = "([0-9a-f]{32})"

func ExtractPageIdsFromPath(path string) []string {
	var pageIds []string
	r := regexp.MustCompile(pageIdPathRegex)
	matches := r.FindAllStringSubmatch(path, -1)
	for _, match := range matches {
		if match != nil && len(match) > 1 {
			pageIds = append(pageIds, match[1])
		}
	}
	return pageIds
}

func ExtractLinkedPageIds(contents string) []string {
	var pageIds []string
	r := regexp.MustCompile(pageIdLinkRegex)
	matches := r.FindAllStringSubmatch(contents, -1)
	for _, match := range matches {
		if match != nil && len(match) > 1 {
			pageIds = append(pageIds, match[1])
		}
	}
	return pageIds
}

func GetContentType(path string) crawler.CrawlContentType {
	fileExtension := filepath.Ext(path)
	if fileExtension == ".md" {
		return crawler.CrawlContentTypeMarkdown
	}

	if fileExtension == ".csv" {
		return crawler.CrawlContentTypeDatabase
	}

	return crawler.CrawlContentTypeUnknown
}
