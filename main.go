package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"spacemen0.github.com/controllers"
	"spacemen0.github.com/helpers"
	"spacemen0.github.com/middlewares"
)

func main() {
	// Initialize the database
	helpers.InitLogger()
	helpers.LoadConfig()
	helpers.InitDB()
	gin.SetMode(gin.ReleaseMode)
	// Set up the Gin router
	router := gin.Default()
	router.HandleMethodNotAllowed = true
	router.Use(cors.Default())

	router.Use(middlewares.LoggerMiddleware())
	// Define the API version group
	v1 := router.Group("/api/v1")

	v1.GET("/search", controllers.Search)
	router.Use(middlewares.DataMiddleware())
	// Define routes for Person
	v1.POST("/people", controllers.CreatePerson)
	v1.GET("/people/:id", controllers.GetPerson)
	v1.PUT("/people/:id", controllers.UpdatePerson)
	v1.DELETE("/people/:id", controllers.DeletePerson)

	// Define routes for Title
	v1.POST("/titles", controllers.CreateTitle)
	v1.GET("/titles/:id", controllers.GetTitle)
	v1.PUT("/titles/:id", controllers.UpdateTitle)
	v1.DELETE("/titles/:id", controllers.DeleteTitle)

	// Start the server
	serverAddress := fmt.Sprintf("%s:%d", helpers.AppConfig.Server.Host, helpers.AppConfig.Server.Port)
	helpers.Log.Printf("Starting server on %s...\n", serverAddress)
	if err := router.Run(serverAddress); err != nil {
		helpers.Log.Fatalf("Failed to start server: %v", err)
	}
}
