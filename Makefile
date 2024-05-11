up-deps:
	docker-compose up qdrant infinity redis neo4j

up-vectordb_sync:
	docker-compose up vectordb_sync

up-notioncrawler:
	docker-compose up --build notioncrawler
