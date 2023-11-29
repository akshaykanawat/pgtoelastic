/*
Version 1.00
Date Created: 2023-11-29
Copyright (c) 2023, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mitchellh/mapstructure"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Hashtags    []struct {
		Name string `json:"name"`
	} `json:"hashtags"`
	Users []struct {
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"users"`
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
				log.Println("Ignored event: %v\n", e)
			}
		}
	}
}

func processNotification(notification Notification, db *sql.DB, client *elasticsearch.TypedClient) {
	switch notification.Table {
	case "user_projects":
		processUserProjectNotification(notification, db, client)
	case "project_hashtags":
		processProjectHashtagNotification(notification, db, client)
	default:
		log.Printf("Unhandled table: %s", notification.Table)
	}
}

func processUserProjectNotification(notification Notification, db *sql.DB, client *elasticsearch.TypedClient) {
	var userProject UserProject
	err := mapstructure.Decode(notification.Data, &userProject)
	if err != nil {
		log.Printf("Error decoding JSON for UserProject: %v", err)
		return
	}

	var user Project
	err = queryProjectData(db, "users", userProject.UserID, &user)
	if err != nil {
		log.Printf("Error querying users table: %v", err)
		return
	}

	var project Project
	err = queryProjectData(db, "projects", userProject.ProjectID, &project)
	if err != nil {
		log.Printf("Error querying projects table: %v", err)
		return
	}

	updateElasticsearchIndex(notification.Operation, client, "projects", fmt.Sprintf("%v", project.ID), project)
}

func processProjectHashtagNotification(notification Notification, db *sql.DB, client *elasticsearch.TypedClient) {
	var projectHashtag ProjectHashtag
	err := mapstructure.Decode(notification.Data, &projectHashtag)
	if err != nil {
		log.Printf("Error decoding JSON for ProjectHashtag: %v", err)
		return
	}
	var project Project
	err = queryProjectData(db, "projects", projectHashtag.ProjectID, &project)
	if err != nil {
		log.Printf("Error querying projects table: %v", err)
		return
	}

	var hashtag Project
	err = queryProjectData(db, "hashtags", projectHashtag.HashtagID, &hashtag)
	if err != nil {
		log.Printf("Error querying hashtags table: %v", err)
		return
	}

	updateElasticsearchIndex(notification.Operation, client, "your-elasticsearch-index", fmt.Sprintf("%v", project.ID), project)
}

func queryProjectData(db *sql.DB, tableName string, id int, project *Project) error {
	return db.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName), id).
		Scan(&project.ID, &project.Name, &project.Slug, &project.Description, &project.CreatedAt)
}

func updateElasticsearchIndex(operation string, client *elasticsearch.TypedClient, indexName, documentID string, project Project) {
	switch operation {
	case "INSERT":
		_, err := client.Index("projects").Id(documentID).Request(project).Do(context.TODO())
		if err != nil {
			log.Printf("Error indexing data into Elasticsearch: %v", err)
		}
	case "DELETE":
		_, err := client.Delete("projects", documentID).Do(context.Background())
		if err != nil {
			log.Printf("Error deleting data from Elasticsearch: %v", err)
		}
	default:
		log.Printf("Unhandled operation: %s", operation)
	}
}
