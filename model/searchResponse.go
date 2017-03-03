package model

type SearchResponse struct {
	TotalResults int64       `json:"total_results"`
	AreaResults  []*Document `json:"area_results,omitempty"`
	Results      []*Document `json:"results"`
}
