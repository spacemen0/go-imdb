package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"spacemen0.github.com/helpers"
	"spacemen0.github.com/models"
)

func TestCreateTitle(t *testing.T) {
	helpers.InitDB() // Initialize the in-memory database

	router := setupRouter()

	// Define a new title
	title := &models.Title{
		ID:            "tt00001",
		PrimaryTitle:  "Example Title",
		OriginalTitle: "Example Original Title",
	}

	// Marshal the title into JSON
	jsonTitle, _ := json.Marshal(title)

	// Create a new POST request
	req, _ := http.NewRequest(http.MethodPost, "/titles", bytes.NewBuffer(jsonTitle))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the response status
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Assert the response body
	var createdTitle models.Title
	err := json.Unmarshal(resp.Body.Bytes(), &createdTitle)
	assert.NoError(t, err)
	assert.Equal(t, "Example Title", createdTitle.PrimaryTitle)
}

func TestGetTitle(t *testing.T) {
	helpers.InitDB() // Initialize the in-memory database

	// First, create a title

	title := &models.Title{
		ID:            "tt00002",
		PrimaryTitle:  "Test Title",
		OriginalTitle: "Test Original Title",
	}
	err := helpers.DB.Create(title).Error
	if err != nil {
		t.Fatalf("Failed to create title: %v", err)
	}

	// Test GET /titles/:id
	router := setupRouter()
	req, _ := http.NewRequest(http.MethodGet, "/titles/tt00002", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	// Assert the response status
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert the response body
	var retrievedTitle models.Title
	err = json.Unmarshal(resp.Body.Bytes(), &retrievedTitle)
	assert.NoError(t, err)
	assert.Equal(t, "Test Title", retrievedTitle.PrimaryTitle)
}

func TestUpdateTitle(t *testing.T) {
	helpers.InitDB() // Initialize the in-memory database

	// First, create a title

	title := &models.Title{
		ID:            "tt00003",
		PrimaryTitle:  "Old Title",
		OriginalTitle: "Old Original Title",
	}
	err := helpers.DB.Create(title).Error
	if err != nil {
		t.Fatalf("Failed to create title: %v", err)
	}

	// Update the title
	updatedTitle := &models.Title{
		ID:            "tt00003",
		PrimaryTitle:  "Updated Title",
		OriginalTitle: "Updated Original Title",
	}
	jsonUpdatedTitle, _ := json.Marshal(updatedTitle)

	req, _ := http.NewRequest(http.MethodPut, "/titles/tt00003", bytes.NewBuffer(jsonUpdatedTitle))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(resp, req)

	// Assert the response status
	assert.Equal(t, http.StatusOK, resp.Code)

	// Assert the response body
	var retrievedTitle models.Title
	err = json.Unmarshal(resp.Body.Bytes(), &retrievedTitle)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", retrievedTitle.PrimaryTitle)
}

func TestDeleteTitle(t *testing.T) {
	helpers.InitDB() // Initialize the in-memory database

	// First, create a title

	title := &models.Title{
		ID:            "tt00004",
		PrimaryTitle:  "Title to Delete",
		OriginalTitle: "Original Title to Delete",
	}
	err := helpers.DB.Create(title).Error
	if err != nil {
		t.Fatalf("Failed to create title: %v", err)
	}

	// Test DELETE /titles/:id
	req, _ := http.NewRequest(http.MethodDelete, "/titles/tt00004", nil)
	resp := httptest.NewRecorder()
	router := setupRouter()
	router.ServeHTTP(resp, req)

	// Assert the response status
	assert.Equal(t, http.StatusNoContent, resp.Code)

	// Verify the title was deleted
	var deletedTitle models.Title
	err = helpers.DB.First(&deletedTitle, "tconst = ?", "tt00004").Error
	assert.Error(t, err) // Expect an error since the title should be deleted
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
