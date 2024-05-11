package controller

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"notioncrawl/services/utils/run_mgr"
	"notioncrawl/services/vector_queue"
)

type ApiController struct {
	neo4j       neo4j.DriverWithContext
	meiliIndex  *meilisearch.Index
	vectorQueue *vector_queue.VectorQueue
	runMgr      *run_mgr.RunMgr
}

func New(neo4j neo4j.DriverWithContext, meiliIndex *meilisearch.Index, vectorQueue *vector_queue.VectorQueue, runMgr *run_mgr.RunMgr) *ApiController {
	return &ApiController{
		neo4j:       neo4j,
		meiliIndex:  meiliIndex,
		vectorQueue: vectorQueue,
		runMgr:      runMgr,
	}
}
