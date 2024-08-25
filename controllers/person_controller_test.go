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

func beforeTesting() *gin.Engine {
	helpers.InitLogger("../server.log")
	helpers.LoadConfig("../config.yaml")
	helpers.InitDB()
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
func afterTesting() {
	helpers.CleanUpDB()
}

func TestCreatePerson(t *testing.T) {
	// Initialize the in-memory database

	router := beforeTesting()

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
	afterTesting()
}

func TestGetPerson(t *testing.T) {
	router := beforeTesting()
	person := &models.Person{
		ID:          "nm00007",
		PrimaryName: "Bob Smith",
	}
	err := helpers.DB.Create(person).Error
	if err != nil {
		t.Fatalf("Failed to create person: %v", err)
	}
	req, _ := http.NewRequest(http.MethodGet, "/people/nm00007", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var retrievedPerson models.Person
	err = json.Unmarshal(resp.Body.Bytes(), &retrievedPerson)
	assert.NoError(t, err)
	assert.Equal(t, "Bob Smith", retrievedPerson.PrimaryName)
	afterTesting()
}

func TestUpdatePerson(t *testing.T) {
	router := beforeTesting()

	person := &models.Person{
		ID:          "nm00008",
		PrimaryName: "Charles Davis",
	}
	err := helpers.DB.Create(person).Error
	assert.NoError(t, err)

	updatedPerson := &models.Person{
		PrimaryName: "Charlie Davis",
	}
	jsonPerson, _ := json.Marshal(updatedPerson)
	req, _ := http.NewRequest(http.MethodPut, "/people/nm00008", bytes.NewBuffer(jsonPerson))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var retrievedPerson models.Person
	err = helpers.DB.First(&retrievedPerson, "nconst = ?", "nm00008").Error
	assert.NoError(t, err)
	assert.Equal(t, "Charlie Davis", retrievedPerson.PrimaryName)
	afterTesting()
}

// Test for deleting a person
func TestDeletePerson(t *testing.T) {
	router := beforeTesting()

	person := &models.Person{
		ID:          "nm00009",
		PrimaryName: "David Johnson",
	}
	err := helpers.DB.Create(person).Error
	assert.NoError(t, err)

	req, _ := http.NewRequest(http.MethodDelete, "/people/nm00009", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)

	var retrievedPerson models.Person
	err = helpers.DB.First(&retrievedPerson, "nconst = ?", "nm00009").Error
	assert.Error(t, err)
	afterTesting()
}
