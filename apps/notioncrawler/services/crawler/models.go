package crawler

type MetaCrawler interface {
	CrawlMeta(blockId string, parentId string) (*CrawledPage, error)
}

type ContentCrawler interface {
	CrawlContent(blockId string) (*ContentCrawlResult, error)
}

type WorkspaceExporter interface {
	Export() error
	GetNextPage() *CrawledPage
	HasNextPage() bool
	Close() error
}

type ContentCrawlResult struct {
	Children []string
	Contents []*CrawlContent
}

type CrawlContentType string

const (
	CrawlContentTypeMarkdown CrawlContentType = "markdown"
	CrawlContentTypeDatabase CrawlContentType = "database"
	CrawlContentTypeUnknown  CrawlContentType = "unknown"
)

type CrawlContent struct {
	Type     CrawlContentType `json:"content_type"`
	FileName string           `json:"file_name"`
	Content  string           `json:"content"`
}

type CrawlNextResult struct {
	CacheMiss bool
}

type CrawlQueueEntry struct {
	PageID   string
	ParentID string
}

type CrawledPage struct {
	PageID      string              `structs:"page_id"`
	ParentID    string              `structs:"-"`
	Url         string              `structs:"url"`
	CrawlResult *ContentCrawlResult `structs:"crawl_result"`
	Hash        int64               `structs:"hash"`
}

type CacheEntry struct {
	PageID     string
	Hash       int64
	ChildPages []string
}

type Options struct {
	ForceUpdateAll bool
	ForceUpdateIds []string
}
