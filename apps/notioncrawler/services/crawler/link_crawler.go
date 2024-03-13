package crawler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"notioncrawl/services/vector_queue"
)

type Neo4jOptions struct {
	Address  string
	Username string
	Password string
}

type Crawler struct {
	id                string
	db                neo4j.DriverWithContext
	vectorQueue       *vector_queue.VectorQueue
	options           *Options
	queue             []*CrawlQueueEntry
	done              []*CrawlQueueEntry
	contentCrawler    ContentCrawler
	metaCrawler       MetaCrawler
	workspaceExporter WorkspaceExporter
}

var defaultOptions = &Options{
	ForceUpdateAll: false,
	ForceUpdateIds: []string{},
}

func New(neo4jOptions Neo4jOptions, vectorQueue *vector_queue.VectorQueue, startPageId string, metaCrawler MetaCrawler, contentCrawler ContentCrawler, workspaceExporter WorkspaceExporter, options *Options) *Crawler {
	if options == nil {
		options = defaultOptions
	}

	driver, err := neo4j.NewDriverWithContext(neo4jOptions.Address, neo4j.BasicAuth(neo4jOptions.Username, neo4jOptions.Password, ""))
	if err != nil {
		panic(err)
	}

	id := uuid.New().String()

	log.Printf("Crawler initiated with id %s", id)

	return &Crawler{
		id:                id,
		db:                driver,
		vectorQueue:       vectorQueue,
		options:           options,
		metaCrawler:       metaCrawler,
		contentCrawler:    contentCrawler,
		workspaceExporter: workspaceExporter,
		queue: []*CrawlQueueEntry{
			{
				PageID:   startPageId,
				ParentID: "",
			},
		},
	}
}

func (s *Crawler) Close() error {
	return s.db.Close(context.Background())
}

func (s *Crawler) QueueSize() int {
	return len(s.queue)
}

func (s *Crawler) HasNext() bool {
	return len(s.queue) > 0
}

func (s *Crawler) dequeue() (*CrawlQueueEntry, error) {
	if len(s.queue) == 0 {
		return nil, errors.New("queue is empty")
	}

	var item *CrawlQueueEntry
	item, s.queue = s.queue[0], s.queue[1:]
	s.done = append(s.done, item)
	return item, nil
}

func (s *Crawler) alreadyEnqueued(entry *CrawlQueueEntry) bool {
	for _, queueEntry := range s.queue {
		if queueEntry.PageID == entry.PageID {
			return true
		}
	}
	return false
}

func (s *Crawler) hasBeenProcessed(entry *CrawlQueueEntry) bool {
	for _, queueEntry := range s.done {
		if queueEntry.PageID == entry.PageID {
			return true
		}
	}
	return false
}

func (s *Crawler) enqueue(entry *CrawlQueueEntry) {
	if s.hasBeenProcessed(entry) || s.alreadyEnqueued(entry) {
		return
	}
	s.queue = append(s.queue, entry)
}

func (s *Crawler) shouldForceUpdate(pageId string) bool {
	if s.options.ForceUpdateAll {
		return true
	}
	for _, id := range s.options.ForceUpdateIds {
		if id == pageId {
			return true
		}
	}
	return false
}

func (s *Crawler) CrawlNext() error {
	queueEntry, err := s.dequeue()
	if err != nil {
		return err
	}

	log.Println(fmt.Sprintf("[info] Id: %s", queueEntry.PageID))

	page, err := s.metaCrawler.CrawlMeta(queueEntry.PageID, queueEntry.ParentID)
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("[info] Url: %s", page.Url))

	var childPageIds []string
	cachedPage := GetCachedPage(s.db, queueEntry.PageID)
	if cachedPage == nil || cachedPage.Hash != page.Hash || s.shouldForceUpdate(queueEntry.PageID) {
		log.Println("[info] Cache MISS")

		content, err := s.contentCrawler.CrawlContent(page.PageID)
		if err != nil {
			return err
		}
		page.CrawlResult = content

		if err := UpdateCache(s.db, page, s.id); err != nil {
			return err
		}

		log.Println("[info] Enqueue for vector db")
		if err := s.vectorQueue.Enqueue(&vector_queue.EnqueuePayload{
			Ids: []string{
				page.PageID,
			},
		}); err != nil {
			fmt.Errorf("[error] Failed to enqueue: %v", err)
		}

		childPageIds = page.CrawlResult.Children
	} else {
		log.Println("[info] Cache HIT")
		childPageIds = cachedPage.ChildPages
	}

	s.done = append(s.done, queueEntry)

	for _, childPageId := range childPageIds {
		s.enqueue(&CrawlQueueEntry{
			ParentID: page.PageID,
			PageID:   childPageId,
		})
	}

	return nil
}

func (s *Crawler) PerformFullBaseExport() error {
	if err := s.workspaceExporter.Export(); err != nil {
		return err
	}

	for s.workspaceExporter.HasNextPage() {
		page := s.workspaceExporter.GetNextPage()
		if page == nil {
			continue
		}
		log.Printf("[info] Page: %s", page.PageID)
		cachedPage := GetCachedPage(s.db, page.PageID)
		if cachedPage == nil || cachedPage.Hash != page.Hash || s.shouldForceUpdate(page.PageID) {
			log.Println("[info] Cache MISS")
			if err := UpdateCache(s.db, page, s.id); err != nil {
				return err
			}
		} else {
			log.Println("[info] Cache HIT")
		}
	}

	return nil
}

func (s *Crawler) Print() {
	print("May print an overview :)")
}
