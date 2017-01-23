package main

import (
	"github.com/ONSdigital/dp-dd-search-api/search"
	"github.com/ONSdigital/dp-publish-pipeline/utils"
	"github.com/ONSdigital/go-ns/log"
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
)

func main() {

	kafkaBrokers := []string{utils.GetEnvironmentVariable("KAFKA_ADDR", "localhost:9092")}
	kafkaConsumerTopic := utils.GetEnvironmentVariable("KAFKA_CONSUMER_TOPIC", "search-index-request")
	kafkaConsumerGroup := utils.GetEnvironmentVariable("KAFKA_CONSUMER_GROUP", "search-index-request")
	elasticSearchNodes := []string{utils.GetEnvironmentVariable("ELASTIC_SEARCH_NODES", "http://127.0.0.1:9200")}
	elasticSearchIndex := utils.GetEnvironmentVariable("ELASTIC_SEARCH_INDEX", "ons")

	log.Debug("Configuration values:", log.Data{
		"KAFKA_ADDR":           kafkaBrokers,
		"KAFKA_CONSUMER_TOPIC": kafkaConsumerTopic,
		"KAFKA_CONSUMER_GROUP": kafkaConsumerGroup,
		"ELASTIC_SEARCH_NODES": elasticSearchNodes,
		"ELASTIC_SEARCH_INDEX": elasticSearchIndex,
	})

	log.Debug("Creating search search client.", nil)
	searchClient, err := search.NewClient(elasticSearchNodes)
	if err != nil {
		log.Error(err, log.Data{"message": "Failed to create Elastic Search client."})
		return
	}
	log.Debug("Elastic Search client Created successfully.", nil)

	log.Debug("Creating Kafka consumer.", nil)
	consumerConfig := cluster.NewConfig()
	kafkaConsumer, err := cluster.NewConsumer(kafkaBrokers, kafkaConsumerGroup, []string{kafkaConsumerTopic}, consumerConfig)
	if err != nil {
		log.Error(err, log.Data{"message": "An error occured creating the Kafka consumer"})
		return
	}
	log.Debug("Kafka consumer created.", nil)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	for {
		select {
		case msg := <-kafkaConsumer.Messages():
			search.ProcessIndexRequest(msg.Value, searchClient, elasticSearchIndex)
		case <-signals:
			log.Debug("Shutting down...", nil)
			err = kafkaConsumer.Close()
			if err != nil {
				log.Error(err, log.Data{"message": "An error occured closing the Kafka consumer"})
			}
			searchClient.Stop()
			log.Debug("Service stopped", nil)
			return
		}
	}
}
