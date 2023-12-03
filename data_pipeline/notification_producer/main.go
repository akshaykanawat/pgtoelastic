/*
Version 1.00
Date Created: 2023-11-28
Copyright (c) 2023, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Notification struct {
	Table     string                 `json:"table"`
	Operation string                 `json:"operation"`
	Data      map[string]interface{} `json:"data"`
}

func main() {
	connStr := PostgresURL //TODO: Fetch from env or aws secrets
	db := openDatabaseConnection(connStr)
	defer db.Close()

	channels := []string{NotificationChannel} // Replace with your channel name
	setupDatabaseListeners(db, channels)

	brokers := BootstrapServer // Replace with your Kafka broker addresses
	producer := setupConfluentKafkaProducer(brokers)
	defer closeConfluentKafkaProducer(producer)

	listener := setupPqListener(connStr)
	defer listener.Close()

	for {
		select {
		case notification := <-listener.Notify:
			handleNotification(notification, producer)
		case <-time.After(90 * time.Second):
			fmt.Println("Received no events for 90 seconds, checking connection")
			go func() {
				listener.Ping()
			}()
		}
	}
}

func openDatabaseConnection(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func setupDatabaseListeners(db *sql.DB, channels []string) {
	for _, channel := range channels {
		_, err := db.Exec("LISTEN " + channel)
		if err != nil {
			panic(err)
		}
	}
}

func setupConfluentKafkaProducer(brokers string) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		panic(err)
	}
	return p
}

func closeConfluentKafkaProducer(producer *kafka.Producer) {
	producer.Close()
}

func setupPqListener(connStr string) *pq.Listener {
	notificationHandler := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			panic(err)
		}
	}

	listener := pq.NewListener(connStr, 10*time.Second, time.Minute, notificationHandler)
	err := listener.Listen(NotificationChannel)
	if err != nil {
		panic(err)
	}
	return listener
}

func handleNotification(notification *pq.Notification, producer *kafka.Producer) {
	fmt.Println("Received notification:", notification.Extra, notification.BePid)
	var dbNotification Notification
	err := json.Unmarshal([]byte(notification.Extra), &dbNotification)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	jsonData, err := json.Marshal(dbNotification)
	if err != nil {
		log.Println("Error converting to JSON:", err)
		return
	}
	kafkaTopic := KafkaTopic
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Key:            []byte(dbNotification.Table),
		Value:          jsonData,
	}

	deliveryChan := make(chan kafka.Event)

	producer.Produce(message, deliveryChan)

	// Wait for delivery report
	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		log.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}
}
