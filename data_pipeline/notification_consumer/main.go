/*
Version 1.00
Date Created: 2023-11-29
Copyright (c) 2023, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	_ "github.com/lib/pq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Notification struct {
	Table     string                 `json:"table"`
	Operation string                 `json:"operation"`
	Data      map[string]interface{} `json:"data"`
}

type UserProject struct {
	ProjectID int `json:"project_id"`
	UserID    int `json:"user_id"`
}

type ProjectHashtag struct {
	HashtagID int `json:"hashtag_id"`
	ProjectID int `json:"project_id"`
}

type User struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CreatedAt  string `json:"created_at"`
	ProjectIds []int  `json:"project_ids"`
}

// Hashtag represents a hashtag entity.
type Hashtag struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	HashtagIds  []int  `json:"hashtag_ids"`
}

func main() {
	// Set up Kafka consumer configuration
	consumerConfig := kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "pgsync-consumer", //TODO: fetch it from config
		"auto.offset.reset": "earliest",
	}

	// Set up Elasticsearch client configuration
	esConfig := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(&consumerConfig)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Subscribe to Kafka topic
	log.Println("topic subscribed")
	err = consumer.SubscribeTopics([]string{"pgsync"}, nil)
	if err != nil {
		log.Fatalf("Error subscribing to Kafka topic: %v", err)
	}

	// Create Elasticsearch client
	esClient, err := elasticsearch.NewTypedClient(esConfig)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %v", err)
	}
	// Create signal channel for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Consume Kafka messages
	run := true
	for run {
		select {
		case sig := <-sigChan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(10000)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				log.Println("kafka_message_received", string(e.Value))
				var notification Notification
				err = json.Unmarshal(e.Value, &notification)
				if err != nil {
					log.Printf("Error decoding JSON: %v", err)
					continue
				}
				processNotification(notification, db, esClient)
			case kafka.Error:
				// Handle Kafka error
				log.Println("kafka_error")
				run = false
				break
			default:
				// Ignore other event types
				log.Printf("Ignored event: %v\n", e)
			}
		}
	}
}

func processNotification(notification Notification, db *sql.DB, client *elasticsearch.TypedClient) {
	switch notification.Table {
	case "user_projects":
		processUserProjectNotification(notification, client)
	case "project_hashtags":
		processProjectHashtagNotification(notification, client)
	case "users":
		processUserNotification(notification, client)
	case "hashtags":
		processHashtagNotification(notification, client)
	case "projects":
		processProjectNotification(notification, client)
	default:
		log.Printf("Unhandled table: %s", notification.Table)
	}
}

func processUserNotification(notification Notification, client *elasticsearch.TypedClient) {
	var user User
	user.ID = int(notification.Data["id"].(float64))
	user.Name = notification.Data["name"].(string)
	user.CreatedAt = notification.Data["created_at"].(string)
	// Update Elasticsearch index
	updateElasticsearchIndex(notification.Operation, client, "users", fmt.Sprintf("%v", user.ID), user)
}

func processHashtagNotification(notification Notification, client *elasticsearch.TypedClient) {
	var hashtag Hashtag
	hashtag.ID = int(notification.Data["id"].(float64))
	hashtag.Name = notification.Data["name"].(string)
	hashtag.CreatedAt = notification.Data["created_at"].(string)

	// Update Elasticsearch index
	updateElasticsearchIndex(notification.Operation, client, "hashtags", fmt.Sprintf("%v", hashtag.ID), hashtag)
}

func processProjectNotification(notification Notification, client *elasticsearch.TypedClient) {
	var project Project
	project.ID = int(notification.Data["id"].(float64))
	project.Name = notification.Data["name"].(string)
	project.Slug = notification.Data["slug"].(string)
	project.Description = notification.Data["description"].(string)
	project.CreatedAt = notification.Data["created_at"].(string)

	// Update Elasticsearch index
	updateElasticsearchIndex(notification.Operation, client, "projects", fmt.Sprintf("%v", project.ID), project)
}

func processUserProjectNotification(notification Notification, client *elasticsearch.TypedClient) {
	var userProject UserProject
	userProject.UserID = int(notification.Data["user_id"].(float64))
	userProject.ProjectID = int(notification.Data["project_id"].(float64))

	operation := notification.Operation
	// Update or delete the Elasticsearch User Index based on the operation
	updateUserProjectIndex(userProject, operation, client)
}
func updateUserProjectIndex(userProject UserProject, operation string, esClient *elasticsearch.TypedClient) {
	// Assuming the Elasticsearch User Index is named "users"
	indexName := "users"

	// Define the update query based on the operation
	var sourceScript string
	switch operation {
	case "INSERT":
		sourceScript = "if (ctx._source.project_ids == null) { ctx._source.project_ids = [] } ctx._source.project_ids.add(params.project_id)"
	case "DELETE":
		sourceScript = "if (ctx._source.containsKey('project_ids')) { ctx._source.project_ids.remove(ctx._source.project_ids.indexOf(params.project_id)) }"
	default:
		log.Printf("Unsupported operation: %s", operation)
		return
	}

	// Define the update query
	query := map[string]interface{}{
		"script": map[string]interface{}{
			"source": sourceScript,
			"lang":   "painless",
			"params": map[string]int{
				"project_id": userProject.ProjectID,
			},
		},
	}

	// Serialize the query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		log.Println("Error marshalling query:", err)
		return
	}

	// Create an io.Reader from the JSON bytes
	bodyReader := bytes.NewBuffer(queryJSON)

	// Define the UpdateRequest
	request := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: fmt.Sprintf("%d", userProject.UserID),
		Body:       bodyReader,
	}

	// Perform the update request
	response, err := request.Do(context.Background(), esClient)
	if err != nil {
		log.Printf("Error updating document. Error: %v", err)
		return
	}

	defer response.Body.Close()

	log.Printf("Document updated successfully. Result: %v", response)
	if response.IsError() {
		var b map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&b); err != nil {
			log.Printf("Error parsing the response body: %v", err)
			return
		}

		log.Printf("Elasticsearch error: %s", b["error"].(map[string]interface{})["reason"].(string))
	}
}

func processProjectHashtagNotification(notification Notification, client *elasticsearch.TypedClient) {
	var projectHashtag ProjectHashtag
	projectHashtag.HashtagID = int(notification.Data["hashtag_id"].(float64))
	projectHashtag.ProjectID = int(notification.Data["project_id"].(float64))
	operation := notification.Operation

	// Update or delete the Elasticsearch Project Index based on the operation
	updateProjectHashtagIndex(projectHashtag, operation, client)
}

func updateProjectHashtagIndex(projectHashtag ProjectHashtag, operation string, esClient *elasticsearch.TypedClient) {
	// Assuming the Elasticsearch Project Index is named "projects"
	indexName := "projects"

	// Define the update query based on the operation
	var sourceScript string
	switch operation {
	case "INSERT":
		sourceScript = "if (ctx._source.hashtag_ids == null){ ctx._source.hashtag_ids = []} ctx._source.hashtag_ids.add(params.hashtag_id)"
	case "DELETE":
		sourceScript = "if (ctx._source.containsKey('hashtag_ids')) { ctx._source.hashtag_ids.remove(ctx._source.hashtag_ids.indexOf(params.hashtag_id)) }"
	default:
		log.Printf("Unsupported operation: %s", operation)
		return
	}

	// Define the update query
	query := map[string]interface{}{
		"script": map[string]interface{}{
			"source": sourceScript,
			"lang":   "painless",
			"params": map[string]int{
				"hashtag_id": projectHashtag.HashtagID,
			},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		log.Println("Error marshalling query:", err)
		return
	}

	// Create an io.Reader from the JSON bytes
	bodyReader := bytes.NewBuffer(queryJSON)

	// Define the UpdateRequest
	request := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: fmt.Sprintf("%d", projectHashtag.ProjectID),
		Body:       bodyReader,
	}

	// Perform the update request
	response, err := request.Do(context.Background(), esClient)
	if err != nil {
		log.Printf("Error updating document. Error: %v", err)
		return
	}

	defer response.Body.Close()

	log.Printf("Document updated successfully. Result: %v", response)
	if response.IsError() {
		var b map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&b); err != nil {
			log.Printf("Error parsing the response body: %v", err)
			return
		}

		log.Printf("Elasticsearch error: %s", b["error"].(map[string]interface{})["reason"].(string))
	}
}

func updateElasticsearchIndex(operation string, client *elasticsearch.TypedClient, indexName, documentID string, data interface{}) {
	switch operation {
	case "INSERT", "UPDATE":
		_, err := client.Index(indexName).Id(documentID).Document(data).Do(context.TODO())
		if err != nil {
			log.Printf("Error indexing data into Elasticsearch: %v", err)
		} else {
			log.Printf("Success: Document %s indexed/updated", documentID)
		}
	case "DELETE":
		_, err := client.Delete("projects", documentID).Do(context.Background())
		if err != nil {
			log.Printf("Error deleting data from Elasticsearch: %v", err)
		} else {
			log.Printf("Success: Document %s deleted", documentID)
		}
	default:
		log.Printf("Unhandled operation: %s", operation)
	}
}
