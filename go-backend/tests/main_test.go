package tests

import (
	"os"
	"testing"

	"github.com/gittuf/visualizer/go-backend/internal/logger"
)

// TestMain sets up the test environment and initializes the logger.
func TestMain(m *testing.M) {
	logger.Initialize()
	defer logger.Sync()

	code := m.Run()

	os.Exit(code)
}
