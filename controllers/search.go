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
		var people []models.Person
		queryStr := "to_tsvector('english', primary_name) @@ plainto_tsquery(?)"
		err = db.Where(queryStr, query).Preload("KnownForTitles").Find(&people).Error
		results = people
	case "title":
		var titles []models.Title
		queryStr := "to_tsvector('english', primary_title || ' ' || original_title) @@ plainto_tsquery(?)"
		err = db.Where(queryStr, query).Preload("Actors").Find(&titles).Error
		results = titles
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}

	c.JSON(http.StatusOK, results)
}
