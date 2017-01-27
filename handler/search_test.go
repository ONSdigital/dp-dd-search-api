package handler_test

import (
	"bytes"
	"errors"
	"github.com/ONSdigital/dp-dd-search-api/handler"
	"github.com/ONSdigital/dp-dd-search-api/model"
	"github.com/ONSdigital/dp-dd-search-api/search/searchtest"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test(t *testing.T) {

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

		Convey("When the query handler is called and the search client returns an error", func() {

			mockSearchClient := searchtest.NewMockSearchClient()
			mockSearchClient.CustomQueryFunc = func(term string) ([]model.Document, error) {
				return nil, errors.New("search client error")
			}
			handler.SearchClient = mockSearchClient
			handler.Search(recorder, request)

			Convey("Then an internal server error code is returned ", func() {
				So(recorder.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}
