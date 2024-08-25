package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"spacemen0.github.com/helpers"
	"spacemen0.github.com/models"
)

// Test for searching persons and titles
func TestSearchController(t *testing.T) {
	router := beforeTesting()

	// Setup: Create test data
	person := &models.Person{
		ID:          "nm00010",
		PrimaryName: "Emily Clark",
	}
	err := helpers.DB.Create(person).Error
	assert.NoError(t, err)

	title := &models.Title{
		ID:            "tt00010",
		PrimaryTitle:  "Test Movie",
		OriginalTitle: "Test Original Movie",
	}
	err = helpers.DB.Create(title).Error
	assert.NoError(t, err)

	// Test search for person
	req, _ := http.NewRequest(http.MethodGet, "/search?query=Emily&by=person", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var personSearchResponse map[string]any
	err = json.Unmarshal(resp.Body.Bytes(), &personSearchResponse)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(personSearchResponse["results"].([]interface{}))) // Expecting one result

	// Test search for title
	req, _ = http.NewRequest(http.MethodGet, "/search?query=Test Movie&by=title", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var titleSearchResponse map[string]any
	err = json.Unmarshal(resp.Body.Bytes(), &titleSearchResponse)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(titleSearchResponse["results"].([]interface{}))) // Expecting one result

	// Test search with invalid query
	req, _ = http.NewRequest(http.MethodGet, "/search?query=NonExisting&by=person", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var emptySearchResponse map[string]any
	err = json.Unmarshal(resp.Body.Bytes(), &emptySearchResponse)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(emptySearchResponse["results"].([]interface{}))) // Expecting zero results
	afterTesting()
}
