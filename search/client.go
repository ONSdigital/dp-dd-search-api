package search

import (
	"fmt"
	"github.com/ONSdigital/dp-dd-search-api/model"
	"gopkg.in/olivere/elastic.v3"
	"reflect"
)

// QueryClient - interface for query functions on the search client.
type QueryClient interface {
	Query(term string, index string) (*model.SearchResponse, error)
	Stop()
}

type elasticSearchClient struct {
	client *elastic.Client
}

// Query - run the given term as a search query
func (elasticSearch *elasticSearchClient) Query(term string, index string) (*model.SearchResponse, error) {

	builder := elasticSearch.client.Search()

	if len(term) > 0 {
		query := elastic.NewQueryStringQuery(term)
		builder.Query(query)
	}

	result, err := builder.
		Index(index).
		From(0).Size(10).
		Pretty(true).
		Do()

	if err != nil {
		return nil, err
	}

	fmt.Printf("Query took %d milliseconds\n", result.TookInMillis)

	var document model.Document
	var documents []*model.Document
	for _, item := range result.Each(reflect.TypeOf(document)) {
		t := item.(model.Document)
		fmt.Printf("Entry %+v\n", t)
		documents = append(documents, &t)
	}

	response := &model.SearchResponse{
		TotalResults: result.TotalHits(),
		Results:      documents,
	}

	return response, nil
}

// Stop the search client
func (elasticSearch *elasticSearchClient) Stop() {
	elasticSearch.client.Stop()
}

// NewClient - Create a new elastic search client instance of QueryClient
func NewClient(nodes []string) (QueryClient, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(nodes...),
		elastic.SetMaxRetries(5))
	if err != nil {
		return nil, err
	}

	return &elasticSearchClient{client}, nil
}
