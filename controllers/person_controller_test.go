package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"spacemen0.github.com/helpers"
	"spacemen0.github.com/models"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/people", CreatePerson)
	r.GET("/people/:id", GetPerson)
	r.PUT("/people/:id", UpdatePerson)
	r.DELETE("/people/:id", DeletePerson)
	r.POST("/titles", CreateTitle)
	r.GET("/titles/:id", GetTitle)
	r.PUT("/titles/:id", UpdateTitle)
	r.DELETE("/titles/:id", DeleteTitle)
	r.GET("/search", Search)
	return r
}

func TestCreatePerson(t *testing.T) {
	helpers.InitDB() // Initialize the in-memory database

	router := setupRouter()

	person := &models.Person{
		ID:          "nm00006",
		PrimaryName: "Alice Johnson",
	}

	jsonPerson, _ := json.Marshal(person)

	req, _ := http.NewRequest(http.MethodPost, "/people", bytes.NewBuffer(jsonPerson))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var createdPerson models.Person
	err := json.Unmarshal(resp.Body.Bytes(), &createdPerson)
	assert.NoError(t, err)
	assert.Equal(t, "Alice Johnson", createdPerson.PrimaryName)
}

func TestGetPerson(t *testing.T) {
	helpers.InitDB() // Initialize the in-memory database

	// First, create a person
	db := helpers.GetDB()
	person := &models.Person{
		ID:          "nm00007",
		PrimaryName: "Bob Smith",
	}
	err := db.Create(person).Error
	if err != nil {
		t.Fatalf("Failed to create person: %v", err)
	}

	// Test GET /people/:id
	router := setupRouter()
	req, _ := http.NewRequest(http.MethodGet, "/people/nm00007", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var retrievedPerson models.Person
	err = json.Unmarshal(resp.Body.Bytes(), &retrievedPerson)
	assert.NoError(t, err)
	assert.Equal(t, "Bob Smith", retrievedPerson.PrimaryName)
}
