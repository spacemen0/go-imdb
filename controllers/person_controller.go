package controllers

import (
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
	db := helpers.GetDB()
	person, err := models.GetPerson(db, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
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
