package model

// Document contains the base properties for all documents stored in Elastic Search.
type Document struct {
	ID   string      `json:"id,omitempty"`
	Type string      `json:"type"`
	Body interface{} `json:"body"`
}
