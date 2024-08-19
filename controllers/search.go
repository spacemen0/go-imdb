package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spacemen0.github.com/helpers"
	"spacemen0.github.com/models"
)

// Search handles GET /search
func Search(c *gin.Context) {
	query := c.Query("query")
	searchType := c.Query("by")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}
	if searchType != "person" && searchType != "title" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity type. Use 'person' or 'title'."})
		return
	}

	db := helpers.GetDB()
	var results any
	var err error

	switch searchType {
	case "person":
		results, err = models.SearchPeople(db, query)
	case "title":
		results, err = models.SearchTitles(db, query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}

	c.JSON(http.StatusOK, results)
}
