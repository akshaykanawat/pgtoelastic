package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"pgsync/server/app/middleware/web"
	"pgsync/server/database"
)

func SearchProjectsByUser(ctx *gin.Context) (*web.JSONResponse, *web.ErrorInterface) {
	// Extract user  from the request or URL parameters
	userId := ctx.Param("id")
	// Query Elasticsearch to get the user details
	user, err := database.GetUserById(userId)
	if err != nil {
		log.Printf("Error querying user details: %v", err)
		return nil, web.NewError(500, "Internal Server Error", "Error querying user details", nil)
	}

	// Check if the user was found
	if user == nil {
		log.Printf("User not found: %s", userId)
		return nil, web.NewHTTPBadRequestError("Not Found", "User not found")
	}

	// Extract user ID and name
	userID := user["id"].(float64)
	userName := user["name"].(string)
	createdAt := user["created_at"].(string)

	// Query Elasticsearch to get the list of project IDs associated with the user
	userProjects, err := database.GetUserProjectsByUserID(int(userID))
	if err != nil {
		log.Printf("Error querying user projects: %v", err)
		return nil, web.NewError(500, "Internal Server Error", "Error querying user projects", nil)
	}

	// Fetch detailed project information using the retrieved project IDs
	projects, err := database.GetProjectsDetails(userProjects)
	if err != nil {
		log.Printf("Error querying project details: %v", err)
		return nil, web.NewError(500, "Internal Server Error", "Error querying project details", nil)
	}

	// Attach hashtag details to each project in the result
	resultProjects, err := attachHashtagsToProjects(projects)
	if err != nil {
		log.Printf("Error attaching hashtag details to projects: %v", err)
		return nil, web.NewError(500, "Internal Server Error", "Error attaching hashtag details to projects", nil)
	}

	// Merge user information with project details
	result := map[string]interface{}{
		"user":     map[string]interface{}{"id": userID, "name": userName, "created_at": createdAt},
		"projects": resultProjects,
	}

	// Return the result as a JSON response
	return web.NewHTTPSuccessResponse(ctx, result), nil
}

func attachHashtagsToProjects(projects []map[string]interface{}) ([]map[string]interface{}, error) {
	// Extract all hashtag IDs from projects
	var hashtagIDs []int
	for _, project := range projects {
		ids, ok := project["hashtag_ids"].([]interface{})
		if !ok {
			continue
		}
		for _, id := range ids {
			if hashtagID, ok := id.(float64); ok {
				hashtagIDs = append(hashtagIDs, int(hashtagID))
			}
		}
	}

	// Query Elasticsearch to get details of all hashtags using the hashtag IDs
	hashtags, err := database.GetHashtagsDetails(hashtagIDs)
	if err != nil {
		log.Printf("Error querying hashtag details: %v", err)
		return nil, err
	}
	// Create a map for faster lookup of hashtag details using ID

	hashtagMap := make(map[int]map[string]interface{})
	for _, hashtag := range hashtags {
		id, ok := hashtag["id"].(float64)
		if ok {
			hashtagMap[int(id)] = hashtag
		}
	}
	log.Printf("hashtag details: %v", hashtagMap)
	// Attach hashtag details to each project in the result
	for _, project := range projects {
		ids, ok := project["hashtag_ids"].([]interface{})
		if !ok {
			continue
		}
		var projectHashtags []map[string]interface{}
		for _, id := range ids {
			if hashtagID, ok := id.(float64); ok {
				if hashtag, exists := hashtagMap[int(hashtagID)]; exists {
					projectHashtags = append(projectHashtags, map[string]interface{}{
						"id":   hashtag["id"],
						"name": hashtag["name"],
					})
				}
			}
		}

		project["hashtags"] = projectHashtags
		delete(project, "hashtag_ids")
		delete(project, "user_ids")
	}
	return projects, nil
}

// SearchProjectsByHashtags handles the search for projects that use specific hashtags.
func SearchProjectsByHashtags(ctx *gin.Context) (*web.JSONResponse, *web.ErrorInterface) {
	// Extract hashtags from the request or URL parameters
	hashtags := ctx.Param("hashtags")
	log.Println(hashtags)

	hashtagDetails, err := database.GetProjectsDetailsByHashtags(hashtags)
	if err != nil {
		log.Printf("Error querying project details: %v", err)
		return nil, web.NewError(500, "Internal Server Error", "Error querying project details", nil)
	}

	// Attach hashtag details and user details to each project in the result
	resultProjects, err := attachDetailsToHashtag(hashtagDetails)
	if err != nil {
		log.Printf("Error attaching details to hashtagDetails: %v", err)
		return nil, web.NewError(500, "Internal Server Error", "Error attaching details to hashtagDetails", nil)
	}

	// Return the result as a JSON response
	return web.NewHTTPSuccessResponse(ctx, map[string]interface{}{"hashtags": resultProjects}), nil
}

// attachDetailsToHashtag attaches hashtag details and user details to each project.
func attachDetailsToHashtag(hashtagDetails []map[string]interface{}) ([]map[string]interface{}, error) {
	// Extract all project IDs and user IDs from hashtagDetails
	var projectIDs []int
	var userIDs []int

	for _, hashtag := range hashtagDetails {
		// Extracting project IDs
		ids, ok := hashtag["project_ids"].([]interface{})
		if ok {
			for _, id := range ids {
				if projectID, ok := id.(float64); ok {
					projectIDs = append(projectIDs, int(projectID))
				}
			}
		}

	}

	// Query Elasticsearch to get details of all projects using the project IDs
	projects, err := database.GetProjectsDetails(projectIDs)
	if err != nil {
		log.Printf("Error querying project details: %v", err)
		return nil, err
	}
	// Create a map for faster lookup of project details using ID
	projectMap := make(map[int]map[string]interface{})
	for _, project := range projects {
		id, ok := project["id"].(float64)
		if ok {
			projectMap[int(id)] = project
		}
		userID, ok := project["user_ids"].([]interface{})
		if ok {
			for _, id := range userID {
				if userID, ok := id.(float64); ok {
					userIDs = append(userIDs, int(userID))
				}
			}
		}
	}
	// Query Elasticsearch to get details of all users using the user IDs
	users, err := database.GetUsersDetails(userIDs)
	if err != nil {
		log.Printf("Error querying user details: %v", err)
		return nil, err
	}
	// Create a map for faster lookup of user details using ID
	userMap := make(map[int]map[string]interface{})
	for _, user := range users {
		id, ok := user["id"].(float64)
		if ok {
			userMap[int(id)] = user
		}
	}
	// Attach project details and user details to each hashtag in the result
	for _, hashtag := range hashtagDetails {
		// Attach project details
		ids, ok := hashtag["project_ids"].([]interface{})
		if ok {
			var projectList []map[string]interface{}
			for _, id := range ids {
				if projectID, ok := id.(float64); ok {
					if project, exists := projectMap[int(projectID)]; exists {

						newUserIDs, ok := project["user_ids"].([]interface{})
						var userList []map[string]interface{}
						if ok {
							for _, eachUser := range newUserIDs {
								if user, exists := userMap[int(eachUser.(float64))]; exists {
									userList = append(userList, map[string]interface{}{
										"id":   user["id"],
										"name": user["name"],
										// Add other user details as needed
									})
								}
								log.Printf("project: %v", project)
							}
						}
						projectList = append(projectList, map[string]interface{}{
							"id":          project["id"],
							"name":        project["name"],
							"slug":        project["slug"],
							"description": project["description"],
							"users":       userList,
							// Add other project details as needed
						})
						log.Printf("project list: %v", projectList)
					}

				}
			}
			hashtag["projects"] = projectList
			delete(hashtag, "project_ids")
		}
		// Attach user details
	}

	return hashtagDetails, nil
}

// FuzzySearchProjects handles the full-text fuzzy search for projects.
func FuzzySearchProjects(ctx *gin.Context) (*web.JSONResponse, *web.ErrorInterface) {
	// Parse the request body to get the search parameters
	var searchParams map[string]string
	if err := ctx.ShouldBindJSON(&searchParams); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return nil, web.NewError(400, "Bad Request", "Error parsing request body", nil)
	}

	// Check if either "slug" or "description" is present in the request body
	var slug, description string
	var ok bool
	if slug, ok = searchParams["slug"]; ok {
		log.Printf("Fuzzy search for slug: %s", slug)
	}

	if description, ok = searchParams["description"]; ok {
		log.Printf("Fuzzy search for description: %s", description)
	}
	response, err := database.FuzzySearchSlugDescription(slug, description)
	if err != nil {
		return nil, web.NewError(400, "Bad Request", "Error parsing request body", nil)
	}
	return web.NewHTTPSuccessResponse(ctx, map[string]interface{}{"projects": response}), nil
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
