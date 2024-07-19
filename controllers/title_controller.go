package controllers

import (
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
	db := helpers.GetDB()
	title, err := models.GetTitle(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Title not found"})
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
