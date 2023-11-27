/*
Version 1.00
Date Created: 2022-06-25
Copyright (c) 2022, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	postgresHost     = "localhost"
	postgresPort     = 5432
	postgresUser     = "postgres"
	postgresPassword = "postgres"
	postgresDB       = "postgres"
)

var (
	userNames    = []string{"John Doe", "Jane Doe", "Alice Smith", "Bob Johnson", "Eva Brown", "Michael Davis"}
	hashtagNames = []string{"Tech", "Travel", "Food", "Fitness", "Science", "Art"}
	projectNames = []string{"Project X", "Awesome App", "Travel Journal", "Healthy Recipes", "AI Research", "Art Gallery"}
)

func main() {
	// Connect to PostgreSQL
	pgConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", postgresHost, postgresPort, postgresUser, postgresPassword, postgresDB)
	pgDB, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer pgDB.Close()

	// Insert seed data maintaining relationships
	err = insertSeedData(pgDB, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Data sync to Elasticsearch complete.")
}

func insertSeedData(pgDB *sql.DB, numRows int) error {
	for i := 1; i <= numRows; i++ {
		// Insert user
		userID, err := insertUser(pgDB, getRandomName(userNames))
		if err != nil {
			return err
		}

		// Insert hashtag
		hashtagID, err := insertHashtag(pgDB, getRandomName(hashtagNames))
		if err != nil {
			return err
		}

		// Insert project
		projectID, err := insertProject(pgDB, getRandomName(projectNames), generateSlug(getRandomName(projectNames)), fmt.Sprintf("Description for %s", getRandomName(projectNames)))
		if err != nil {
			return err
		}

		// Link project with hashtag
		err = linkProjectHashtag(pgDB, projectID, hashtagID)
		if err != nil {
			return err
		}

		// Link user with project
		err = linkUserProject(pgDB, userID, projectID)
		if err != nil {
			return err
		}
	}
	return nil
}

func getRandomName(names []string) string {
	return names[rand.Intn(len(names))]
}

func insertUser(pgDB *sql.DB, name string) (int, error) {
	var userID int
	err := pgDB.QueryRow("INSERT INTO users (name, created_at) VALUES ($1, $2) RETURNING id", name, time.Now()).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func insertHashtag(pgDB *sql.DB, name string) (int, error) {
	var hashtagID int
	err := pgDB.QueryRow("INSERT INTO hashtags (name, created_at) VALUES ($1, $2) RETURNING id", name, time.Now()).Scan(&hashtagID)
	if err != nil {
		return 0, err
	}
	return hashtagID, nil
}

func insertProject(pgDB *sql.DB, name, slug, description string) (int, error) {
	var projectID int
	err := pgDB.QueryRow("INSERT INTO projects (name, slug, description, created_at) VALUES ($1, $2, $3, $4) RETURNING id", name, slug, description, time.Now()).Scan(&projectID)
	if err != nil {
		return 0, err
	}
	return projectID, nil
}

func linkProjectHashtag(pgDB *sql.DB, projectID, hashtagID int) error {
	_, err := pgDB.Exec("INSERT INTO project_hashtags (hashtag_id, project_id) VALUES ($1, $2)", hashtagID, projectID)
	return err
}

func linkUserProject(pgDB *sql.DB, userID, projectID int) error {
	_, err := pgDB.Exec("INSERT INTO user_projects (project_id, user_id) VALUES ($1, $2)", projectID, userID)
	return err
}

func generateSlug(name string) string {
	// Generate a lowercase slug by replacing spaces with hyphens
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}