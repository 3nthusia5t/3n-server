package sqlite

import (
	"os"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

func TestInit(t *testing.T) {
	// Define your test directory and file path
	filePath := "../../rsrc/test_database/database.db" // Change this to your test directory

	// Clean up previous test files if they exist
	os.Remove(filePath)
	defer os.Remove(filePath)

	// Call your Init function
	dbManager := Init(filePath)
	if dbManager == nil {
		t.Errorf("Failed to initialize database")
	}
	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("Init did not create file: %s", filePath)
	}
}
