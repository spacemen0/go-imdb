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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := helpers.GetDB()
	if err := models.CreatePerson(db, &person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, person)
}

// GetPerson handles GET /persons/:id
func GetPerson(c *gin.Context) {
	id := c.Param("id")
	verbose := c.DefaultQuery("verbose", "false")
	db := helpers.GetDB()
	var person *models.Person
	var err error

	if verbose == "true" {
		// If verbose=true, preload all associated data
		person, err = models.GetPerson(db, id, true)
	} else {
		// If verbose=false, only preload limited associations (e.g., IDs)
		person, err = models.GetPerson(db, id, false)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve person"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	person.ID = id
	db := helpers.GetDB()
	if err := models.UpdatePerson(db, &person); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}
	c.JSON(http.StatusOK, person)
}

// DeletePerson handles DELETE /persons/:id
func DeletePerson(c *gin.Context) {
	id := c.Param("id")
	db := helpers.GetDB()
	if err := models.DeletePerson(db, id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
