package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"pgsync/server/models"
)

var esClient *elasticsearch.Client

func getUserDetails(userID string) (*models.User, error) {
	// Elasticsearch query for user details
	res, err := esClient.Get("users", userID)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	var user models.User
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Function to query Elasticsearch for projects created by a particular user
func getProjectsCreatedByUser(userID string) ([]models.Project, error) {
	// Elasticsearch query for projects created by a particular user
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"user_id": userID,
			},
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.Search(esClient.Search.WithIndex("projects"), esClient.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	var searchResult models.ESSearchResult
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		return nil, err
	}

	var projectsDetails []models.Project
	for _, hit := range searchResult.Hits.Hits {
		var project models.Project
		err := json.Unmarshal(hit.Source, &project)
		if err != nil {
			return nil, err
		}
		projectsDetails = append(projectsDetails, project)
	}

	return projectsDetails, nil
}

// Function to query Elasticsearch for projects that use specific hashtags
func getProjectsByHashtags(hashtags []string) ([]models.Project, error) {
	// Elasticsearch query for projects that use specific hashtags
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"hashtags": hashtags,
			},
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := esClient.Search(esClient.Search.WithIndex("projects"), esClient.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	var searchResult models.ESSearchResult
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		return nil, err
	}

	var projectsDetails []models.Project
	for _, hit := range searchResult.Hits.Hits {
		var project models.Project
		err := json.Unmarshal(hit.Source, &project)
		if err != nil {
			return nil, err
		}
		projectsDetails = append(projectsDetails, project)
	}

	return projectsDetails, nil
}

// Function to perform full-text fuzzy search for projects
func fullTextFuzzySearch(query string) ([]models.Project, error) {
	// Elasticsearch query for full-text fuzzy search
	fuzzyQuery := map[string]interface{}{
		"fuzzy": map[string]interface{}{
			"slug": map[string]interface{}{
				"value":     query,
				"fuzziness": "AUTO",
			},
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(fuzzyQuery); err != nil {
		return nil, err
	}

	res, err := esClient.Search(esClient.Search.WithIndex("projects"), esClient.Search.WithBody(&buf))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	var searchResult models.ESSearchResult
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		return nil, err
	}

	var projectsDetails []models.Project
	for _, hit := range searchResult.Hits.Hits {
		var project models.Project
		err := json.Unmarshal(hit.Source, &project)
		if err != nil {
			return nil, err
		}
		projectsDetails = append(projectsDetails, project)
	}

	return projectsDetails, nil
}

// Function to query Elasticsearch for hashtags used by the projects along with details of the users
func getHashtagsAndUsersForProjects(projectIDs []string) ([]string, []string, error) {
	// Elasticsearch query to get hashtags and users for projects
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"ids": map[string]interface{}{
				"values": projectIDs,
			},
		},
		"_source": map[string]interface{}{
			"includes": []string{"hashtags", "user.*"},
		},
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, nil, err
	}
	res, err := esClient.Search(esClient.Search.WithIndex("projects"), esClient.Search.WithBody(&buf))
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, nil, fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	var searchResult models.ESSearchResult
	err = json.NewDecoder(res.Body).Decode(&searchResult)
	if err != nil {
		return nil, nil, err
	}

	var hashtags []string
	var users []string

	for _, hit := range searchResult.Hits.Hits {
		var project models.Project
		err := json.Unmarshal(hit.Source, &project)
		if err != nil {
			return nil, nil, err
		}

		hashtags = append(hashtags, project.Hashtags...)
		users = append(users, project.Users...)
	}

	return hashtags, users, nil
}

// ConnectEsClient function to initialize Elasticsearch client
func ConnectEsClient() (*elasticsearch.Client, error) {
	var esConfig = elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %v", err)
		return nil, err
	}

	return client, nil
}
