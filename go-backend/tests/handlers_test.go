package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gittuf/visualizer/go-backend/internal/handlers"
	"github.com/gittuf/visualizer/go-backend/internal/models"
	"github.com/gittuf/visualizer/go-backend/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// setupRouter initializes the Gin engine and registers the API routes for testing.
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Register all routes here
	r.POST("/commits", handlers.ListCommits)
	r.POST("/metadata", handlers.GetMetadata)
	r.POST("/commits-local", handlers.ListCommitsLocal)
	r.POST("/metadata-local", handlers.GetMetadataLocal)

	return r
}

// TestListCommits_Success verifies that the /commits endpoint correctly returns a list of commits from a remote repository.
func TestListCommits_Success(t *testing.T) {
	// Setup a "remote" repo
	remotePath, _, cleanupRemote := helpers.SetupTestRepo(t)
	defer cleanupRemote()

	r := setupRouter()

	reqBody := models.CommitsRequest{
		URL: remotePath, // Use local path as URL
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/commits", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var commits []models.Commit
	err := json.Unmarshal(w.Body.Bytes(), &commits)
	assert.NoError(t, err)
	assert.NotEmpty(t, commits)
}

// TestGetMetadata_Success verifies that the /metadata endpoint correctly returns the requested metadata file from a remote repository.
func TestGetMetadata_Success(t *testing.T) {
	// Setup a "remote" repo
	remotePath, commitHash, cleanupRemote := helpers.SetupTestRepo(t)
	defer cleanupRemote()

	r := setupRouter()

	reqBody := models.MetadataRequest{
		URL:    remotePath,
		Commit: commitHash,
		File:   "root.json",
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/metadata", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var metadata map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &metadata)
	assert.NoError(t, err)
	assert.Equal(t, "root", metadata["type"])
}

// TestListCommitsLocal_Success verifies that the /commits-local endpoint correctly returns a list of commits from a local repository.
func TestListCommitsLocal_Success(t *testing.T) {
	// Setup a repo
	repoPath, _, cleanup := helpers.SetupTestRepo(t)
	defer cleanup()

	r := setupRouter()

	reqBody := models.CommitsLocalRequest{
		Path: repoPath,
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/commits-local", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var commits []models.Commit
	err := json.Unmarshal(w.Body.Bytes(), &commits)
	assert.NoError(t, err)
	assert.NotEmpty(t, commits)
}

// TestGetMetadataLocal_Success verifies that the /metadata-local endpoint correctly returns the requested metadata file from a local repository.
func TestGetMetadataLocal_Success(t *testing.T) {
	// Setup a repo
	repoPath, commitHash, cleanup := helpers.SetupTestRepo(t)
	defer cleanup()

	r := setupRouter()

	reqBody := models.MetadataLocalRequest{
		Path:   repoPath,
		Commit: commitHash,
		File:   "root.json",
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/metadata-local", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var metadata map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &metadata)
	assert.NoError(t, err)
	assert.Equal(t, "root", metadata["type"])
}

// TestListCommitsLocal_InvalidPath verifies that the /commits-local endpoint returns a 400 Bad Request for an invalid repository path.
func TestListCommitsLocal_InvalidPath(t *testing.T) {
	r := setupRouter()

	reqBody := models.CommitsLocalRequest{
		Path: "/invalid/path",
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/commits-local", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetMetadataLocal_InvalidPath verifies that the /metadata-local endpoint returns a 400 Bad Request for an invalid repository path.
func TestGetMetadataLocal_InvalidPath(t *testing.T) {
	r := setupRouter()

	reqBody := models.MetadataLocalRequest{
		Path:   "/invalid/path",
		Commit: "HEAD",
		File:   "root.json",
	}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/metadata-local", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestListCommits_MissingURL verifies that the /commits endpoint returns a 400 Bad Request when the URL is missing.
func TestListCommits_MissingURL(t *testing.T) {
	r := setupRouter()

	reqBody := models.CommitsRequest{}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/commits", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestGetMetadata_MissingURL verifies that the /metadata endpoint returns a 400 Bad Request when the URL is missing.
func TestGetMetadata_MissingURL(t *testing.T) {
	r := setupRouter()

	reqBody := models.MetadataRequest{}
	jsonValue, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/metadata", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
