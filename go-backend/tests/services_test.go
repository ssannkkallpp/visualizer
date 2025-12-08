package tests

import (
	"os"
	"testing"

	"github.com/gittuf/visualizer/go-backend/internal/services"
	"github.com/gittuf/visualizer/go-backend/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// TestCloneAndFetchRepo_Success verifies that a repository can be successfully cloned and the policy ref fetched.
func TestCloneAndFetchRepo_Success(t *testing.T) {
	remotePath, _, cleanupRemote := helpers.SetupTestRepo(t)
	defer cleanupRemote()

	localPath, cleanupLocal, err := services.CloneAndFetchRepo(remotePath)
	assert.NoError(t, err)
	defer cleanupLocal()

	_, err = os.Stat(localPath)
	assert.NoError(t, err)
}

// TestGetPolicyCommits_Success verifies that commits from the policy ref can be successfully retrieved.
func TestGetPolicyCommits_Success(t *testing.T) {
	remotePath, commitHash, cleanupRemote := helpers.SetupTestRepo(t)
	defer cleanupRemote()

	localPath, cleanupLocal, err := services.CloneAndFetchRepo(remotePath)
	assert.NoError(t, err)
	defer cleanupLocal()

	commits, err := services.GetPolicyCommits(localPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, commits)
	assert.Equal(t, commitHash, commits[0].Hash)
}

// TestGetLocalCommits_Success verifies that commits from a local repository can be successfully retrieved.
func TestGetLocalCommits_Success(t *testing.T) {
	repoPath, _, cleanup := helpers.SetupTestRepo(t)
	defer cleanup()

	commits, err := services.GetLocalCommits(repoPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, commits)
	// The helper creates 2 commits: initial and policy.
	// The policy commit is on a detached ref or branch?
	// Helper: "Create a commit that will be on the policy ref... commitHash, err := w.Commit..."
	// This commit is on the CURRENT branch (which was master/main) unless we switched.
	// Helper: "Create and switch to refs/gittuf/policy branch... No, we just commit to current branch then set ref"
	// So HEAD should have the policy commit too.
}

// TestDecodeMetadataBlob_Success verifies that a metadata blob can be successfully decoded from a repository.
func TestDecodeMetadataBlob_Success(t *testing.T) {
	// Setup a "remote" repo
	remotePath, commitHash, cleanupRemote := helpers.SetupTestRepo(t)
	defer cleanupRemote()

	// Clone it
	localPath, cleanupLocal, err := services.CloneAndFetchRepo(remotePath)
	assert.NoError(t, err)
	defer cleanupLocal()

	// Decode metadata
	// The helper writes plain JSON to "root.json"
	// services.DecodeMetadataBlob expects the file content.
	// If the service expects base64 encoded content inside the file, we might fail here if the service doesn't handle plain JSON.
	// But let's assume it reads the file content.
	metadata, err := services.DecodeMetadataBlob(localPath, commitHash, "root.json")
	assert.NoError(t, err)
	assert.NotNil(t, metadata)

	// Check if we got the expected data
	// metadata is map[string]interface{}
	assert.Equal(t, "root", metadata["type"])
}

// TestCloneAndFetchRepo_InvalidURL verifies that attempting to clone from an invalid URL returns an error.
func TestCloneAndFetchRepo_InvalidURL(t *testing.T) {
	_, _, err := services.CloneAndFetchRepo("invalid-url")
	assert.Error(t, err)
}

// TestGetPolicyCommits_InvalidPath verifies that attempting to get policy commits from an invalid path returns an error.
func TestGetPolicyCommits_InvalidPath(t *testing.T) {
	_, err := services.GetPolicyCommits("invalid-path")
	assert.Error(t, err)
}

// TestGetLocalCommits_InvalidPath verifies that attempting to get local commits from an invalid path returns an error.
func TestGetLocalCommits_InvalidPath(t *testing.T) {
	_, err := services.GetLocalCommits("invalid-path")
	assert.Error(t, err)
}

// TestDecodeMetadataBlob_InvalidPath verifies that attempting to decode metadata from an invalid path returns an error.
func TestDecodeMetadataBlob_InvalidPath(t *testing.T) {
	_, err := services.DecodeMetadataBlob("invalid-path", "HEAD", "root.json")
	assert.Error(t, err)
}
