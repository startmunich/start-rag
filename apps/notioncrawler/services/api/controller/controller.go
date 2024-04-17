package controller

import (
	"github.com/meilisearch/meilisearch-go"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ApiController struct {
	neo4j      neo4j.DriverWithContext
	meiliIndex *meilisearch.Index
}

func New(neo4j neo4j.DriverWithContext, meiliIndex *meilisearch.Index) *ApiController {
	return &ApiController{
		neo4j:      neo4j,
		meiliIndex: meiliIndex,
	}
}
