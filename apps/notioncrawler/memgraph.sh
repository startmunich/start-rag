#!/bin/bash
docker run --name memgraph -p 7687:7687 -p 7444:7444 -p 3000:3000 memgraph/memgraph-platform