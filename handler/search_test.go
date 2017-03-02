package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/ONSdigital/dp-dd-search-api/config"
	"github.com/ONSdigital/dp-dd-search-api/handler"
	"github.com/ONSdigital/dp-dd-search-api/model"
	"github.com/ONSdigital/dp-dd-search-api/search/searchtest"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {

	var documents []*model.Document
	documents = append(documents, &model.Document{
		ID:   "document 1",
		Type: "dataset",
	})
	documents = append(documents, &model.Document{
		ID:   "document 2",
		Type: "dataset",
	})

	var areas []*model.Document
	areas = append(areas, &model.Document{
		ID:   "areaID",
		Type: "area",
	})

	mockSearchResponse := &model.SearchResponse{
		Results:      documents,
		AreaResults:  areas,
		TotalResults: 2,
	}

	mockSuggestResponse := &model.SearchResponse{
		Results:      documents,
		TotalResults: 2,
	}

	Convey("Given a search request", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte("{not a valid document}"))
		request, _ := http.NewRequest("GET", "/search?q=armed", requestBodyReader)

		Convey("When the query handler is called", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			handler.SearchClient = mockSearchClient
			handler.Search(recorder, request)

			Convey("Then the search client is called with the expected parameters", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)
				So(mockSearchClient.QueryRequests[0], ShouldEqual, "armed")
			})
		})
	})

	Convey("Given a search request", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte("{not a valid document}"))
		request, _ := http.NewRequest("GET", "/search?q=armed", requestBodyReader)

		Convey("When the query handler is called", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			mockSearchClient.CustomQueryFunc = func(term string, index string) (*model.SearchResponse, error) {
				return mockSearchResponse, nil
			}
			handler.SearchClient = mockSearchClient
			handler.Search(recorder, request)

			Convey("Then the response contains the content returned from the search client.", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)

				actualResponse := &model.SearchResponse{}
				_ = json.Unmarshal(recorder.Body.Bytes(), actualResponse)

				So(actualResponse.TotalResults, ShouldEqual, mockSearchResponse.TotalResults)
				So(actualResponse.Results[0].ID, ShouldEqual, mockSearchResponse.Results[0].ID)
				So(actualResponse.Results[1].ID, ShouldEqual, mockSearchResponse.Results[1].ID)
				So(actualResponse.AreaResults[0].ID, ShouldEqual, mockSearchResponse.AreaResults[0].ID)
			})
		})
	})

	Convey("Given a search request", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte("{not a valid document}"))
		request, _ := http.NewRequest("GET", "/search?q=armed", requestBodyReader)

		Convey("When the query handler is called and the search client returns an error", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			mockSearchClient.CustomQueryFunc = func(term string, index string) (*model.SearchResponse, error) {
				return nil, errors.New("search client error")
			}
			handler.SearchClient = mockSearchClient
			handler.Search(recorder, request)

			Convey("Then the result has an empty results array", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)

				actualResponse := &model.SearchResponse{}
				_ = json.Unmarshal(recorder.Body.Bytes(), actualResponse)

				So(actualResponse.TotalResults, ShouldEqual, 0)
			})
		})
	})

	Convey("Given a search request", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte("{not a valid document}"))
		request, _ := http.NewRequest("GET", "/search?q=armed", requestBodyReader)

		Convey("When the query handler is called and the search area client returns an error", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			mockSearchClient.CustomQueryFunc = func(term string, index string) (*model.SearchResponse, error) {

				if index == config.ElasticSearchIndex {
					return mockSearchResponse, nil
				}

				return nil, errors.New("search client error")
			}
			handler.SearchClient = mockSearchClient
			handler.Search(recorder, request)

			Convey("Then the result has an empty results array", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)

				actualResponse := &model.SearchResponse{}
				_ = json.Unmarshal(recorder.Body.Bytes(), actualResponse)
				So(len(actualResponse.AreaResults), ShouldEqual, 0)
			})
		})
	})


	Convey("Given a suggest request", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte("{not a valid document}"))
		request, _ := http.NewRequest("GET", "/search?q=armed", requestBodyReader)

		Convey("When the query handler is called", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			handler.SearchClient = mockSearchClient
			handler.Suggest(recorder, request)

			Convey("Then the search client is called with the expected parameters", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)
				So(mockSearchClient.QueryRequests[0], ShouldEqual, "armed")
			})
		})
	})

	Convey("Given a suggest request", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte("{not a valid document}"))
		request, _ := http.NewRequest("GET", "/suggest?q=armed", requestBodyReader)

		Convey("When the query handler is called", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			mockSearchClient.CustomQueryFunc = func(term string, index string) (*model.SearchResponse, error) {
				return mockSuggestResponse, nil
			}
			handler.SearchClient = mockSearchClient
			handler.Suggest(recorder, request)

			Convey("Then the response contains the content returned from the search client.", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)

				actualResponse := &model.SearchResponse{}
				_ = json.Unmarshal(recorder.Body.Bytes(), actualResponse)

				So(actualResponse.TotalResults, ShouldEqual, mockSuggestResponse.TotalResults)
				So(actualResponse.Results[0].ID, ShouldEqual, mockSuggestResponse.Results[0].ID)
			})
		})
	})

	Convey("Given a suggest request", t, func() {

		recorder := httptest.NewRecorder()
		requestBodyReader := bytes.NewReader([]byte("{not a valid document}"))
		request, _ := http.NewRequest("GET", "/suggst?q=armed", requestBodyReader)

		Convey("When the query handler is called and the search client returns an error", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			mockSearchClient.CustomQueryFunc = func(term string, index string) (*model.SearchResponse, error) {
				return nil, errors.New("search client error")
			}
			handler.SearchClient = mockSearchClient
			handler.Suggest(recorder, request)

			Convey("Then the result has an empty results array", func() {
				So(recorder.Code, ShouldEqual, http.StatusOK)

				actualResponse := &model.SearchResponse{}
				_ = json.Unmarshal(recorder.Body.Bytes(), actualResponse)

				So(actualResponse.TotalResults, ShouldEqual, 0)
			})
		})
	})
}
