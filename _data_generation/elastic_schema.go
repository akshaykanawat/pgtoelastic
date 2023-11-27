/*
Version 1.00
Date Created: 2022-06-25
Copyright (c) 2022, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"net/http"
)

const (
	elasticURL   = "http://your-elasticsearch-host:9200"
	indexName    = "projects"
	indexMapping = `
	{
		"mappings": {
			"properties": {
				"name": {
					"type": "text"
				},
				"slug": {
					"type": "keyword"
				},
				"description": {
					"type": "text"
				},
				"hashtags": {
					"type": "keyword"
				},
				"user": {
					"type": "text"
				},
				"created_at": {
					"type": "date",
					"format": "strict_date_optional_time||epoch_millis"
				}
			}
		}
	}`
)

func main() {
	// Create Elasticsearch client
	cfg := elasticsearch.Config{
		CloudID: "9120838b227147f3aac6a7cc80bb8b2c:dXMtY2VudHJhbDEuZ2NwLmNsb3VkLmVzLmlvJDc2ZTRkZGY4NzI3ODRhNjJhMGVhNmQ5Y2VkZmM3NGZkJDMxMDQ0ZjQyMzQ1NzQyNGJhOWIzZWIzMTMxNTExNDNm",
		APIKey:  "ZVhiekJvd0JyVmxRUG1OQm9Dd1U6NUE2SVlvaXhRNmFIMzU1UWx0VkRNUQ==",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Create the index with the specified mapping
	//infores, err := es.Info()
	//if err != nil {
	//	log.Fatalf("Error getting response: %s", err)
	//}
	//	fmt.Println(infores)
	err = createIndex(es)
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}

	fmt.Println("Index created successfully.")
}

func createIndex(es *elasticsearch.Client) error {
	// Check if the index already exists
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == http.StatusNotFound {
			// Index doesn't exist, create it
			return createIndexWithMapping(es)
		}
		return fmt.Errorf("Error checking if index exists: %s", res.Status())
	}

	fmt.Printf("Index '%s' already exists.\n", indexName)
	return nil
}

func createIndexWithMapping(es *elasticsearch.Client) error {
	// Create the index with the specified mapping
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  bytes.NewReader([]byte(indexMapping)),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Error creating index: %s", res.Status())
	}

	return nil
}
