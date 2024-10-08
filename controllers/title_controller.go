package controllers

import (
	"errors"
	"gorm.io/gorm"
	"net/http"

	"spacemen0.github.com/helpers"
	"spacemen0.github.com/models"

	"github.com/gin-gonic/gin"
)

// CreateTitle handles POST /titles
func CreateTitle(c *gin.Context) {
	var title models.Title
	if err := c.ShouldBindJSON(&title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}
	if err := models.CreateTitle(helpers.DB, &title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create title", "details": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, title)
}

// GetTitle handles GET /titles/:id
func GetTitle(c *gin.Context) {
	id := c.Param("id")
	verbose := c.DefaultQuery("verbose", "false")

	var title *models.Title
	var err error

	if verbose == "true" {
		title, err = models.GetTitle(helpers.DB, id, true)
	} else {
		title, err = models.GetTitle(helpers.DB, id, false)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Title not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve title", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, title)
}

// UpdateTitle handles PUT /titles/:id
func UpdateTitle(c *gin.Context) {
	id := c.Param("id")
	var title models.Title
	if err := c.ShouldBindJSON(&title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	title.ID = id

	if err := models.UpdateTitle(helpers.DB, &title); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Title not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update title", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, title)
}

// DeleteTitle handles DELETE /titles/:id
func DeleteTitle(c *gin.Context) {
	id := c.Param("id")

	if err := models.DeleteTitle(helpers.DB, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Title not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete title", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
