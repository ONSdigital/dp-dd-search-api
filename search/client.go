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
	Suggest(term string) (*model.SearchResponse, error)
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
		documents = append(documents, &t)
	}

	response := &model.SearchResponse{
		TotalResults: result.TotalHits(),
		Results:      documents,
	}

	return response, nil
}

func (elasticSearch *elasticSearchClient) Suggest(term string) (*model.SearchResponse, error) {
	builder := elasticSearch.client.Search()

	if len(term) <= 0 {
		response := &model.SearchResponse{
			TotalResults: 0,
			Results:      []*model.Document{},
		}
		return response, nil
	}

	fetchSourceContext := elastic.NewFetchSourceContext(true).Exclude("body.dimensions", "body.metadata")
	query := elastic.NewMatchPhrasePrefixQuery("body.title", term)
	builder.FetchSourceContext(fetchSourceContext).Query(query)


	result, err := builder.
		Index("dd").
		Size(10).
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

	//logger := log.New(os.Stdout, "", log.LstdFlags)
	client, err := elastic.NewClient(
		elastic.SetURL(nodes...),
		elastic.SetMaxRetries(5))

	// Add these lines into the client initialisation to enable query logging.
	//elastic.SetInfoLog(logger)
	//elastic.SetTraceLog(logger)

	if err != nil {
		return nil, err
	}

	return &elasticSearchClient{client}, nil
}
