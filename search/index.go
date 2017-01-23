package search

import (
	"encoding/json"
	"github.com/ONSdigital/go-ns/log"
)

type IndexRequest struct {
	Type string            `json:"type"`
	ID   string            `json:"id"`
	Data map[string]string `json:"data"`
}

func ProcessIndexRequest(msg []byte, elasticSearchClient Client, elasticSearchIndex string) {

	var request IndexRequest
	err := json.Unmarshal(msg, &request)
	if err != nil {
		log.Debug("Failed to parse json request data", nil)
		return
	}

	log.Debug("Indexing document", log.Data{
		"Document": request,
	})

	err = elasticSearchClient.Index(elasticSearchIndex, request.Type, request.ID, request.Data)

	if err != nil {
		log.Error(err, nil)
	}
}
