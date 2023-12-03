package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"pgsync/server/config"
	"strings"
	"sync"
)

var (
	esClient     *elasticsearch.Client
	esClientOnce sync.Once
)

func GetUserProjectsByUserID(userID int) ([]int, error) {
	_, err := ConnectEsClient()
	if err != nil {
		return nil, err
	}
	indexName := config.UserIndex
	// Define the query to get user projects
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"id": userID,
			},
		},
	}

	// Serialize the query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Create an io.Reader from the JSON bytes
	bodyReader := strings.NewReader(string(queryJSON))

	// Define the SearchRequest
	request := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bodyReader,
	}

	// Perform the search request
	response, err := request.Do(context.Background(), esClient)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract project IDs from the response
	var projectIDs []int
	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	for _, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if !ok {
			continue
		}
		projectIDsField, ok := source["project_ids"].([]interface{})
		if !ok {
			continue
		}
		for _, projectID := range projectIDsField {
			if id, ok := projectID.(float64); ok {
				projectIDs = append(projectIDs, int(id))
			}
		}
	}

	return projectIDs, nil
}

// getUserByName queries Elasticsearch to get the user details by name.
func GetUserById(userID string) (map[string]interface{}, error) {
	_, err := ConnectEsClient()
	if err != nil {
		return nil, err
	}
	indexName := config.UserIndex
	// Define the query to get user details by name
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"id": userID,
			},
		},
	}
	// Serialize the query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Create an io.Reader from the JSON bytes
	bodyReader := strings.NewReader(string(queryJSON))

	// Define the SearchRequest
	request := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bodyReader,
	}

	// Perform the search request
	response, err := request.Do(context.Background(), esClient)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract user details from the response
	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok || len(hits) == 0 {
		return nil, nil // User not found
	}

	source, ok := hits[0].(map[string]interface{})["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return source, nil
}

// getProjectsDetails queries Elasticsearch to fetch detailed information about projects using their IDs.
func GetProjectsDetails(projectIDs []int) ([]map[string]interface{}, error) {
	_, err := ConnectEsClient()
	if err != nil {
		return nil, err
	}
	indexName := config.ProjectIndex

	// Define the query to get project details by project IDs
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"id": projectIDs,
			},
		},
	}

	// Serialize the query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	// Create an io.Reader from the JSON bytes
	bodyReader := strings.NewReader(string(queryJSON))

	// Define the SearchRequest
	request := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bodyReader,
	}

	// Perform the search request
	response, err := request.Do(context.Background(), esClient)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Extract project details from the response
	var projects []map[string]interface{}
	hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}
	for _, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if !ok {
			continue
		}
		projects = append(projects, source)
	}

	return projects, nil
}

func GetHashtagsDetails(hashtagIDs []int) ([]map[string]interface{}, error) {
	// Filter out null values from hashtagIDs
	_, err := ConnectEsClient()
	if err != nil {
		return nil, err
	}
	filteredHashtagIDs := make([]int, 0, len(hashtagIDs))
	for _, id := range hashtagIDs {
		if id != 0 { // Assuming 0 represents a null value, adjust this condition if needed
			filteredHashtagIDs = append(filteredHashtagIDs, id)
		}
	}

	// Check if all values were null
	if len(filteredHashtagIDs) == 0 {
		return nil, nil
	}

	// Build the terms query for hashtag IDs
	termsQuery := map[string]interface{}{
		"terms": map[string]interface{}{
			"id": filteredHashtagIDs,
		},
	}

	// Create the search request
	searchRequest := map[string]interface{}{
		"query": termsQuery,
	}

	// Convert the search request to JSON
	searchJSON, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, err
	}

	// Execute the search query against Elasticsearch
	res, err := esClient.Search(
		esClient.Search.WithIndex(config.HashtagIndex),
		esClient.Search.WithBody(bytes.NewReader(searchJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check if the response is successful
	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	// Parse the response to get the hashtag details
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Extract the hits from the response
	hits, ok := response["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid Elasticsearch response format")
	}
	// Extract the source of each hit (hashtag details)
	var hashtagDetails []map[string]interface{}
	for _, hit := range hits {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if ok {
			hashtagDetails = append(hashtagDetails, source)
		}
	}

	return hashtagDetails, nil
}

// GetProjectsDetailsByHashtags fetches project details based on specified hashtags.
func GetProjectsDetailsByHashtags(hashtag string) ([]map[string]interface{}, error) {
	_, err := ConnectEsClient()
	if err != nil {
		return nil, err
	}
	indexName := config.HashtagIndex

	// Build the terms query for hashtag
	termQuery := map[string]interface{}{
		"match": map[string]interface{}{
			"name": hashtag,
		},
	}

	// Create the search request
	searchRequest := map[string]interface{}{
		"query": termQuery,
	}

	// Convert the search request to JSON
	searchJSON, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, err
	}

	// Execute the search query against Elasticsearch
	res, err := esClient.Search(
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(bytes.NewReader(searchJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check if the response is successful
	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	// Parse the response to get the project details
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}
	// Extract the hits from the response
	hits, ok := response["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid Elasticsearch response format")
	}
	// Extract the source of each hit (project details)
	var projectDetails []map[string]interface{}
	for _, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if ok {
			projectDetails = append(projectDetails, source)
		}
	}

	return projectDetails, nil
}

func GetUsersDetails(userIDs []int) ([]map[string]interface{}, error) {
	_, err := ConnectEsClient()
	if err != nil {
		return nil, err
	}
	indexName := config.UserIndex

	// Build the terms query for user IDs
	termsQuery := map[string]interface{}{
		"terms": map[string]interface{}{
			"id": userIDs,
		},
	}

	// Create the search request
	searchRequest := map[string]interface{}{
		"query": termsQuery,
	}

	// Convert the search request to JSON
	searchJSON, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, err
	}

	// Execute the search query against Elasticsearch
	res, err := esClient.Search(
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(bytes.NewReader(searchJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check if the response is successful
	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	// Parse the response to get the user details
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Extract the hits from the response
	hits, ok := response["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid Elasticsearch response format")
	}

	// Extract the source of each hit (user details)
	var userDetails []map[string]interface{}
	for _, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if ok {
			userDetails = append(userDetails, source)
		}
	}

	return userDetails, nil
}

func FuzzySearchSlugDescription(slug string, description string) ([]map[string]interface{}, error) {
	// Build the fuzzy search query for slug
	_, err := ConnectEsClient()
	if err != nil {
		return nil, err
	}
	log.Printf("Values: %v, %v", description, slug)
	slugQuery := buildFuzzyQuery("slug.keyword", slug)

	// Build the fuzzy search query for description
	descriptionQuery := buildFuzzyQuery("description.keyword", description)

	// Combine the queries using a boolean OR
	combinedQuery := map[string]interface{}{
		"bool": map[string]interface{}{
			"should": []interface{}{slugQuery, descriptionQuery},
		},
	}

	// Create the search request
	searchRequest := map[string]interface{}{
		"query": combinedQuery,
	}
	log.Printf("query: %v", searchRequest)
	// Convert the search request to JSON
	searchJSON, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, err
	}

	// Execute the search query against Elasticsearch
	res, err := esClient.Search(
		esClient.Search.WithIndex(config.ProjectIndex),
		esClient.Search.WithBody(bytes.NewReader(searchJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check if the response is successful
	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	// Parse the response to get the project details
	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	// Extract the hits from the response
	hits, ok := response["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid Elasticsearch response format")
	}

	// Extract the source of each hit (project details)
	var projectDetails []map[string]interface{}
	for _, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if ok {
			projectDetails = append(projectDetails, source)
		}
	}

	return projectDetails, nil
}

func buildFuzzyQuery(field string, value string) map[string]interface{} {
	return map[string]interface{}{
		"fuzzy": map[string]interface{}{
			field: map[string]interface{}{
				"value":     value,
				"fuzziness": "AUTO", // You can adjust the fuzziness level as needed
			},
		},
	}
}

// ConnectEsClient returns the cached Elasticsearch client or creates a new one if it doesn't exist.
func ConnectEsClient() (*elasticsearch.Client, error) {
	esClientOnce.Do(func() {
		var esConfig = elasticsearch.Config{
			Addresses: []string{"http://localhost:9200"},
		}

		client, err := elasticsearch.NewClient(esConfig)
		if err != nil {
			log.Fatalf("Error creating Elasticsearch client: %v", err)
			return
		}

		esClient = client
	})

	return esClient, nil
}
