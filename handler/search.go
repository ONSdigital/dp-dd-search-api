package handler

import (
	"encoding/json"
	"github.com/ONSdigital/dp-dd-search-api/search"
	"github.com/ONSdigital/go-ns/log"
	"net/http"
)

// SearchClient - the dependency used to interact with elastic search.
var SearchClient search.QueryClient

// Search - HTTP handler for accepting search query requests.
func Search(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query().Get("q")

	results, err := SearchClient.Query(query, "dd")
	if err != nil {
		log.Error(err, log.Data{
			"message": "Error running a search query.",
			"query":   query})
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	areaResults, err := SearchClient.Query(query, "areas")
	if err != nil {
		log.Error(err, log.Data{
			"message": "Error running a search query.",
			"query":   query})
		w.WriteHeader(http.StatusInternalServerError)
		return
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
