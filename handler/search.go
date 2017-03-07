package handler

import (
	"encoding/json"
	"github.com/ONSdigital/dp-dd-search-api/model"
	"github.com/ONSdigital/dp-dd-search-api/search"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
)

// SearchClient - the dependency used to interact with elastic search.
var SearchClient search.QueryClient

// Search - HTTP handler for accepting search query requests.
func Search(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query().Get("q")
	filter := req.URL.Query().Get("filter")

	results, err := SearchClient.Query(query, filter, "dd")
	if err != nil {
		log.Error(err, log.Data{
			"message": "Error running a search query.",
			"query":   query})
		results = &model.SearchResponse{}
		results.Results = make([]*model.Document, 0)
	}

	areaResults, err := SearchClient.Query(query, "", "areas")
	if err != nil {
		log.Error(err, log.Data{
			"message": "Error running a search query.",
			"query":   query})
		results.AreaResults = make([]*model.Document, 0)
	}

	if areaResults != nil {
		results.AreaResults = areaResults.Results
	}

	responseJSON, err := json.Marshal(results)
	if err != nil {
		log.Error(err, log.Data{
			"message": "Error serialising the search results to JSON",
			"query":   query,
			"results": results})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Error(err, log.Data{"message": "Error writing response body"})
	}
}

func Suggest(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("q")

	results, err := SearchClient.Suggest(query)
	if err != nil {
		log.Error(err, log.Data{
			"message": "Error running an auto-complete (suggest) query.",
			"query":   query})
		results = &model.SearchResponse{}
		results.Results = make([]*model.Document, 0)
	}

	responseJSON, err := json.Marshal(results)
	if err != nil {
		log.Error(err, log.Data{
			"message": "Error serialising the search results to JSON",
			"query":   query,
			"results": results})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	_, err = w.Write(responseJSON)
	if err != nil {
		log.Error(err, log.Data{"message": "Error writing response body"})
	}
}
