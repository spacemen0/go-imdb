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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := helpers.GetDB()
	if err := models.CreateTitle(db, &title); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, title)
}

// GetTitle handles GET /titles/:id
func GetTitle(c *gin.Context) {
	id := c.Param("id")
	verbose := c.DefaultQuery("verbose", "false")
	db := helpers.GetDB()
	var title *models.Title
	var err error
	if verbose == "true" {
		title, err = models.GetTitle(db, id, true)
	} else if verbose == "false" {
		title, err = models.GetTitle(db, id, false)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Title not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve title"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	title.ID = id
	db := helpers.GetDB()
	if err := models.UpdateTitle(db, &title); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Title not found"})
		return
	}
	c.JSON(http.StatusOK, title)
}

// DeleteTitle handles DELETE /titles/:id
func DeleteTitle(c *gin.Context) {
	id := c.Param("id")
	db := helpers.GetDB()
	if err := models.DeleteTitle(db, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Title not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
