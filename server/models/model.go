/*
Version 1.00
Date Created: 2022-06-25
Copyright (c) 2022, Akshay Singh Kanawat
Author: Akshay Singh Kanawat
*/
package models

import "encoding/json"

type User struct {
	ID         string   `json:"id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	ProjectIDs []string `json:"project_ids"`
}

type Project struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	Hashtags    []string `json:"hashtags"`
	Users       []string `json:"users"`
}

type ESSearchResult struct {
	Took     int64     `json:"took"`
	TimedOut bool      `json:"timed_out"`
	Shards   ShardInfo `json:"_shards"`
	Hits     HitsInfo  `json:"hits"`
}

type ShardInfo struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type HitsInfo struct {
	Total    TotalInfo `json:"total"`
	MaxScore float64   `json:"max_score"`
	Hits     []HitInfo `json:"hits"`
}

type TotalInfo struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type HitInfo struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}
