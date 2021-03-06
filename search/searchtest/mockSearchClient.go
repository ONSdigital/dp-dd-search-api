package searchtest

import (
	"github.com/ONSdigital/dp-dd-search-api/model"
	"github.com/ONSdigital/dp-dd-search-api/search"
)

// Checks the MockSearchClient satisfies the IndexingClient interface
var _ search.QueryClient = (*MockSearchClient)(nil)

// NewMockSearchClient creates a new instance of MockSearchClient
func NewMockSearchClient() *MockSearchClient {
	return &MockSearchClient{}
}

// MockSearchClient provides a mock implementation of QueryClient
type MockSearchClient struct {
	QueryRequests   []string
	FilterRequests  []string
	CustomQueryFunc func(term string, index string) (*model.SearchResponse, error)
}

// Query - just capture the query request for later assertion.
func (elasticSearch *MockSearchClient) Query(term string, filter string, index string) (*model.SearchResponse, error) {

	if elasticSearch.CustomQueryFunc != nil {
		return elasticSearch.CustomQueryFunc(term, index)
	}

	elasticSearch.QueryRequests = append(elasticSearch.QueryRequests, term)
	elasticSearch.FilterRequests = append(elasticSearch.FilterRequests, filter)
	return nil, nil
}

// Query - just capture the query request for later assertion.
func (elasticSearch *MockSearchClient) Suggest(term string) (*model.SearchResponse, error) {

	if elasticSearch.CustomQueryFunc != nil {
		return elasticSearch.CustomQueryFunc(term, "dd")
	}

	elasticSearch.QueryRequests = append(elasticSearch.QueryRequests, term)
	return nil, nil
}

// Stop - mock implementation does nothing.
func (elasticSearch *MockSearchClient) Stop() {}
