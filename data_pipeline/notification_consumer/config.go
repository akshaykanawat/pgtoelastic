/*
Version 1.00
Date Created: 2022-06-25
Copyright (c) 2022, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package main

const PostgresURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
const KafkaTopic = "pgsync"
const BootstrapServer = "localhost:9092"
const ConsumerGroup = "pgsync-consumer"
const ElastisearchURL = "http://localhost:9200"
const TableUserProjects = "user_projects"
const TableProjectHashtags = "project_hashtags"
const TableUsers = "users"
const TableHashtags = "hashtags"
const TableProjects = "projects"
const IndexUsers = "users"
const IndexHashtags = "hashtags"
const IndexProjects = "projects"
const OperationInsert = "INSERT"
const OperationUpdate = "UPDATE"
const OperationDelete = "DELETE"
