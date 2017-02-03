package main

import (
	"github.com/ONSdigital/dp-dd-search-api/config"
	"github.com/ONSdigital/dp-dd-search-api/handler"
	"github.com/ONSdigital/dp-dd-search-api/search"
	"github.com/ONSdigital/go-ns/handlers/healthcheck"
	"github.com/ONSdigital/go-ns/handlers/requestID"
	"github.com/ONSdigital/go-ns/handlers/timeout"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gorilla/pat"
	"github.com/justinas/alice"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	config.Load()

	log.Debug("Creating search search client.", nil)
	searchClient, err := search.NewClient(config.ElasticSearchNodes, config.ElasticSearchIndex)
	if err != nil {
		log.Error(err, log.Data{"message": "Failed to create Elastic Search client."})
		os.Exit(1)
	}

	handler.SearchClient = searchClient

	exitCh := make(chan struct{})

	listenForHTTPRequests(exitCh)
	waitForInterrupt(searchClient, exitCh)
}

func listenForHTTPRequests(exitCh chan struct{}) {

	go func() {
		router := pat.New()
		router.Get("/healthcheck", healthcheck.Handler)
		router.Get("/search", handler.Search)
		log.Debug("Starting HTTP server", log.Data{"bind_addr": config.BindAddr})

		middleware := []alice.Constructor{
			requestID.Handler(16),
			log.Handler,
			corsHandler,
			timeout.Handler(10 * time.Second),
		}
		alice := alice.New(middleware...).Then(router)

		server := &http.Server{
			Addr:         config.BindAddr,
			Handler:      alice,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		if err := server.ListenAndServe(); err != nil {
			log.Error(err, nil)
		}

		log.Debug("HTTP server has stopped.", nil)
		exitCh <- struct{}{}
	}()
}

func corsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		h.ServeHTTP(w, req)
	})
}

func waitForInterrupt(searchClient search.QueryClient, exitCh chan struct{}) {

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)

	select {
	case <-signals:
		log.Debug("OS signal receieved.", nil)
		shutdown(searchClient)
	case <-exitCh:
		log.Debug("Notification received on exit channel.", nil)
		shutdown(searchClient)
	}
}

func shutdown(searchClient search.QueryClient) {
	log.Debug("Shutting down.", nil)
	searchClient.Stop()
	log.Debug("Service stopped", nil)
}
