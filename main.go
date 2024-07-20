package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"spacemen0.github.com/controllers"
	"spacemen0.github.com/helpers"
)

func main() {
	// Initialize the database
	helpers.InitDB()
	gin.SetMode(gin.ReleaseMode)
	// Set up the Gin router
	router := gin.Default()

	// Define the API version group
	v1 := router.Group("/api/v1")

	// Define routes for Person
	v1.POST("/people", controllers.CreatePerson)
	v1.GET("/people/:id", controllers.GetPerson)
	v1.PUT("/people/:id", controllers.UpdatePerson)
	v1.DELETE("/people/:id", controllers.DeletePerson)

	// Define routes for Title if needed
	v1.POST("/titles", controllers.CreateTitle)
	v1.GET("/titles/:id", controllers.GetTitle)
	v1.PUT("/titles/:id", controllers.UpdateTitle)
	v1.DELETE("/titles/:id", controllers.DeleteTitle)

	// Start the server
	log.Println("Starting server on :8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
