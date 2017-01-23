package searchtest

func NewMockSearchClient() *MockSearchClient {
	return &MockSearchClient{}
}

type MockSearchClient struct {
	IndexRequests []IndexRequest
}

type IndexRequest struct {
	Index        string
	DocumentType string
	Id           string
	Body         interface{}
}

func (elasticSearch *MockSearchClient) Index(index string, documentType string, id string, body interface{}) error {
	elasticSearch.IndexRequests = append(elasticSearch.IndexRequests, IndexRequest{
		Index:        index,
		DocumentType: documentType,
		Id:           id,
		Body:         body,
	})
	return nil
}

func (elasticSearch *MockSearchClient) Stop() {}
