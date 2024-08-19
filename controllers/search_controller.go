package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spacemen0.github.com/helpers"
	"spacemen0.github.com/models"
	"strconv"
)

// Search handles GET /search
func Search(c *gin.Context) {
	query := c.Query("query")
	searchType := c.Query("by")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}
	if searchType != "person" && searchType != "title" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity type. Use 'person' or 'title'."})
		return
	}

	offset := (page - 1) * limit

	db := helpers.GetDB()
	var results any
	var err error

	switch searchType {
	case "person":
		results, err = models.SearchPeople(db, query, limit, offset)
	case "title":
		results, err = models.SearchTitles(db, query, limit, offset)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
		return
	}

	// Return paginated results with metadata
	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"limit":   limit,
		"results": results,
	})
}
