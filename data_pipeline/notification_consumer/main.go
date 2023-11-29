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
	_ "github.com/lib/pq"
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

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// Hashtag represents a hashtag entity.
type Hashtag struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Project represents a project entity.
type ProjectDB struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Hashtags    []Hashtag `json:"hashtags"`
	Users       []User    `json:"users"`
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
	userProject.UserID = int(notification.Data["user_id"].(float64))
	userProject.ProjectID = int(notification.Data["project_id"].(float64))
	var project Project
	err := queryUserProjectData(db, userProject.UserID, &project)
	if err != nil {
		log.Printf("Error querying users table: %v", err)
		return
	}
	updateElasticsearchIndex(notification.Operation, client, "projects", fmt.Sprintf("%v", project.ID), project)
}

func processProjectHashtagNotification(notification Notification, db *sql.DB, client *elasticsearch.TypedClient) {
	var projectHashtag ProjectHashtag
	projectHashtag.HashtagID = int(notification.Data["hashtag_id"].(float64))
	projectHashtag.ProjectID = int(notification.Data["project_id"].(float64))
	err := mapstructure.Decode(notification.Data, &projectHashtag)
	var project Project
	err = queryProjectHashtagData(db, projectHashtag.ProjectID, &project)
	if err != nil {
		log.Printf("Error querying projects table: %v", err)
		return
	}
	updateElasticsearchIndex(notification.Operation, client, "projects", fmt.Sprintf("%v", project.ID), project)
}

func queryProjectHashtagData(db *sql.DB, id int, project *Project) error {
	var projectScanArgs = []interface{}{&project.ID, &project.Name, &project.Slug, &project.Description, &project.CreatedAt}

	// Query projects table
	projectQuery := "SELECT * FROM projects WHERE id = $1"
	err := db.QueryRow(projectQuery, id).Scan(projectScanArgs...)
	if err != nil {
		return err
	}

	// Query project_hashtags table
	hashtagsQuery := "SELECT h.* FROM hashtags h INNER JOIN project_hashtags ph ON h.id = ph.hashtag_id WHERE ph.project_id = $1"
	rows, err := db.Query(hashtagsQuery, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	var hashtags []Hashtag
	for rows.Next() {
		hashtag := Hashtag{}
		err := rows.Scan(&hashtag.ID, &hashtag.Name, &hashtag.CreatedAt)
		if err != nil {
			return err
		}
		hashtags = append(hashtags, hashtag)
	}
	fmt.Println("hash: ", hashtags)
	project.Hashtags = hashtags

	return nil
}

func queryUserProjectData(db *sql.DB, id int, project *Project) error {
	var projectScanArgs = []interface{}{&project.ID, &project.Name, &project.Slug, &project.Description, &project.CreatedAt}

	// Query projects table
	projectQuery := "SELECT * FROM projects WHERE id = $1"
	err := db.QueryRow(projectQuery, id).Scan(projectScanArgs...)
	if err != nil {
		return err
	}

	// Query user_projects table
	usersQuery := "SELECT u.* FROM users u INNER JOIN user_projects up ON u.id = up.user_id WHERE up.project_id = $1"
	rows, err := db.Query(usersQuery, id)
	if err != nil {
		return err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt)
		if err != nil {
			return err
		}
		users = append(users, user)
	}
	project.Users = users

	return nil
}

func updateElasticsearchIndex(operation string, client *elasticsearch.TypedClient, indexName, documentID string, project Project) {
	switch operation {
	case "INSERT":
		_, err := client.Index(indexName).Id(documentID).Request(project).Do(context.TODO())
		if err != nil {
			log.Printf("Error indexing data into Elasticsearch: %v", err)
		} else {
			log.Printf("success inserted")
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
