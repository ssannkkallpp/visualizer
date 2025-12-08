package helpers

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// SetupTestRepo creates a temporary git repository with a commit on the gittuf policy branch
func SetupTestRepo(t *testing.T) (string, string, func()) {
	tempDir, err := os.MkdirTemp("", "gittuf-viz-test-repo-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to init git repo: %v", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to get worktree: %v", err)
	}

	// Create a dummy file and commit to main
	dummyFile := filepath.Join(tempDir, "README.md")
	if err := os.WriteFile(dummyFile, []byte("# Test Repo"), 0644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to write dummy file: %v", err)
	}
	if _, err := w.Add("README.md"); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to add file: %v", err)
	}
	if _, err := w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	}); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to commit: %v", err)
	}

	// Create a commit that will be on the policy ref
	rootJSON := `{"type":"root", "expires":"2030-01-01T00:00:00Z"}`
	rootB64 := base64.StdEncoding.EncodeToString([]byte(rootJSON))
	envelope := fmt.Sprintf(`{"payload": "%s", "signatures": []}`, rootB64)

	// Create metadata directory
	metadataDir := filepath.Join(tempDir, "metadata")
	if err := os.MkdirAll(metadataDir, 0755); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create metadata dir: %v", err)
	}

	policyFile := filepath.Join(metadataDir, "root.json")
	if err := os.WriteFile(policyFile, []byte(envelope), 0644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to write policy file: %v", err)
	}
	if _, err := w.Add("metadata/root.json"); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to add policy file: %v", err)
	}

	commitHash, err := w.Commit("Add root.json", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Gittuf Admin",
			Email: "admin@gittuf.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to commit policy: %v", err)
	}

	// Create the specific ref
	refName := "refs/gittuf/policy"

	// Using plumbing to set the ref
	ref := plumbing.NewHashReference(plumbing.ReferenceName(refName), commitHash)
	if err := repo.Storer.SetReference(ref); err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to set policy ref: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, commitHash.String(), cleanup
}

// CreateBase64Metadata encodes the given content string into a base64 string.
func CreateBase64Metadata(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}
