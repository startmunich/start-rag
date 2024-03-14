package controller

import "github.com/neo4j/neo4j-go-driver/v5/neo4j"

type ApiController struct {
	neo4j neo4j.DriverWithContext
}

func New(neo4j neo4j.DriverWithContext) *ApiController {
	return &ApiController{
		neo4j: neo4j,
	}
}
