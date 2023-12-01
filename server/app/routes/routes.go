package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"pgsync/server/app/handlers"
	"pgsync/server/app/middleware"
)

const GetProjectsByUsers = "/user/:userID"
const GetProjectByHashtag = "/hashtags/:hashtags"
const SearchProject = "/search"
const SendError = "/error"
const HealthCheck = "/health-check"

func SetupServer() *gin.Engine {
	log.Printf("setting_up_routes...")
	r := gin.Default()
	pgCore := r.Group("/v1/projects")
	addV1Routes(pgCore)
	addProjectRoutes(pgCore) // Add new project routes
	return r
}

// addV1Routes: add routes to v1 endpoint
func addV1Routes(rg *gin.RouterGroup) {
	ping(rg)
	sendError(rg)
	getUser(rg)
	getHashtag(rg)
	searchProject(rg)
}

func addProjectRoutes(rg *gin.RouterGroup) {
	// Add new project routes
	rg.GET(GetProjectsByUsers, middleware.ServeEndpoint(handlers.SearchProjectsByUser))
	rg.GET(GetProjectByHashtag, middleware.ServeEndpoint(handlers.SearchProjectsByHashtags))
	rg.GET(SearchProject, middleware.ServeEndpoint(handlers.FuzzySearchProjects))
}

// Function to add ping route
func ping(rg *gin.RouterGroup) {
	rg.GET(HealthCheck, middleware.ServeEndpoint(handlers.SendPing))
}

// Function to add error route
func sendError(rg *gin.RouterGroup) {
	rg.GET(SendError, middleware.ServeEndpoint(handlers.SendError))
}

func getUser(rg *gin.RouterGroup) {
	rg.GET(SendError, middleware.ServeEndpoint(handlers.SendError))
}

func getHashtag(rg *gin.RouterGroup) {
	rg.GET(SendError, middleware.ServeEndpoint(handlers.SendError))
}

func searchProject(rg *gin.RouterGroup) {
	rg.GET(SendError, middleware.ServeEndpoint(handlers.SendError))
}
