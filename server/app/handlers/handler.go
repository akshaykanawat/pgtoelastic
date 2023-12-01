package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"pgsync/server/app/middleware/web"
)

func SearchProjectsByUser(ctx *gin.Context) (*web.JSONResponse, *web.ErrorInterface) {
	// Extract user ID from the request or URL parameters
	userID := ctx.Param("userID")
	log.Println(userID)
	// Use userID to query Elasticsearch for projects created by the user
	// Implement Elasticsearch query logic here

	// Return the result as a JSON response
	return web.NewHTTPSuccessResponse(ctx, map[string]interface{}{}), nil
}

// SearchProjectsByHashtags handles the search for projects that use specific hashtags.
func SearchProjectsByHashtags(ctx *gin.Context) (*web.JSONResponse, *web.ErrorInterface) {
	// Extract hashtags from the request or URL parameters
	hashtags := ctx.Param("hashtags")
	log.Println(hashtags)
	// Use hashtags to query Elasticsearch for projects
	// Implement Elasticsearch query logic here

	// Return the result as a JSON response
	return web.NewHTTPSuccessResponse(ctx, map[string]interface{}{}), nil
}

// FuzzySearchProjects handles the full-text fuzzy search for projects.
func FuzzySearchProjects(ctx *gin.Context) (*web.JSONResponse, *web.ErrorInterface) {
	// Extract search query from the request or URL parameters
	searchQuery := ctx.Param("searchQuery")
	log.Println(searchQuery)
	// Use searchQuery to perform fuzzy search on project slug and description in Elasticsearch
	// Implement Elasticsearch query logic here

	// Return the result as a JSON response
	return web.NewHTTPSuccessResponse(ctx, map[string]interface{}{}), nil
}

func SendPing(ctx *gin.Context) (*web.JSONResponse, *web.ErrorInterface) {
	log.Print("Health check call")
	message := map[string]interface{}{
		"message": "pong",
	}
	response := web.NewHTTPSuccessResponse(ctx, message)
	//response.JSON()
	return response, nil
}

func SendError(ctx *gin.Context) (*web.JSONResponse, *web.ErrorInterface) {
	return nil, web.NewHTTPBadRequestError("Downtime", map[string]interface{}{
		"Error": "Temporary Downtime",
	})
}
