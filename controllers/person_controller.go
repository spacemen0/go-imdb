package controllers

import (
	"errors"
	"gorm.io/gorm"
	"net/http"

	"spacemen0.github.com/helpers"
	"spacemen0.github.com/models"

	"github.com/gin-gonic/gin"
)

// CreatePerson handles POST /persons
func CreatePerson(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	if err := models.CreatePerson(helpers.DB, &person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create person", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, person)
}

// GetPerson handles GET /persons/:id
func GetPerson(c *gin.Context) {
	id := c.Param("id")
	verbose := c.DefaultQuery("verbose", "false")
	var person *models.Person
	var err error

	if verbose == "true" {
		person, err = models.GetPerson(helpers.DB, id, true)
	} else {
		person, err = models.GetPerson(helpers.DB, id, false)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve person", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, person)
}

// UpdatePerson handles PUT /persons/:id
func UpdatePerson(c *gin.Context) {
	id := c.Param("id")
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data", "details": err.Error()})
		return
	}

	person.ID = id
	if err := models.UpdatePerson(helpers.DB, &person); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, person)
}

// DeletePerson handles DELETE /persons/:id
func DeletePerson(c *gin.Context) {
	id := c.Param("id")
	if err := models.DeletePerson(helpers.DB, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
