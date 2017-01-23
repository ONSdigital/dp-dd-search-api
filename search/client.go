package search

import (
	"gopkg.in/olivere/elastic.v3"
)

type Client interface {
	Index(index string, documentType string, id string, body interface{}) error
	Stop()
}

type ElasticSearchClient struct {
	client *elastic.Client
}

func NewClient(nodes []string) (Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(nodes...),
		elastic.SetMaxRetries(5))
	if err != nil {
		return nil, err
	}

	return &ElasticSearchClient{client}, nil
}

func (elasticSearch *ElasticSearchClient) Index(index string, documentType string, id string, body interface{}) error {
	_, err := elasticSearch.client.Index().
		Index(index).
		Type(documentType).
		Id(id).
		BodyJson(body).
		Refresh(true).
		Do()

	return err
}

func (elasticSearch *ElasticSearchClient) Stop() {
	elasticSearch.client.Stop()
}
