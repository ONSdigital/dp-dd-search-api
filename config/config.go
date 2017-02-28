package config

import (
	"github.com/ONSdigital/go-ns/log"
	"os"
	"strings"
)

// BindAddr - the port to listen for HTTP requests
var BindAddr string

// ElasticSearchNodes - a list of elastic search nodes to use
var ElasticSearchNodes []string

// ElasticSearchIndex - the Elastic Search index name to write to.
var ElasticSearchIndex string

// Load any defined environment variables into the configuration.
func Load() {
	BindAddr = getEnvironmentVariable("BIND_ADDR", ":20051")
	ElasticSearchNodes = strings.Split(getEnvironmentVariable("ELASTIC_SEARCH_NODES", "http://localhost:9200"), ",")
	ElasticSearchIndex = getEnvironmentVariable("ELASTIC_SEARCH_INDEX", "dd")

	log.Debug("Loaded configuration values:", log.Data{
		"BIND_ADDR":            BindAddr,
		"ELASTIC_SEARCH_NODES": ElasticSearchNodes,
		"ELASTIC_SEARCH_INDEX": ElasticSearchIndex,
	})
}

func getEnvironmentVariable(name string, defaultValue string) string {
	environmentValue := os.Getenv(name)
	if environmentValue != "" {
		return environmentValue
	}

	return defaultValue
}
