package search_test

import (
	"encoding/json"
	"github.com/ONSdigital/dp-dd-search-api/search"
	"github.com/ONSdigital/dp-dd-search-api/search/searchtest"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestProcessor(t *testing.T) {

	Convey("Given a new index request", t, func() {
		expectedRequest := search.IndexRequest{
			Data: map[string]string{
				"got data?": "yes I do",
			},
			ID:   "123",
			Type: "thetype",
		}
		indexRequestString, _ := json.Marshal(expectedRequest)
		searchClient := searchtest.NewMockSearchClient()
		searchIndex := "searchindex"

		Convey("When the index request is processed", func() {
			search.ProcessIndexRequest(indexRequestString, searchClient, searchIndex)

			Convey("Then the search client is called with the expected parameters", func() {
				var actualRequest searchtest.IndexRequest = searchClient.IndexRequests[0]

				So(actualRequest.DocumentType, ShouldEqual, expectedRequest.Type)
				So(actualRequest.Id, ShouldEqual, expectedRequest.ID)

				requestDataIsEqual := reflect.DeepEqual(actualRequest.Body, expectedRequest.Data)
				So(requestDataIsEqual, ShouldBeTrue)
			})
		})
	})
}
