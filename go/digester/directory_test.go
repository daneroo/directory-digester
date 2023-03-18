package digester

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDirectory(t *testing.T) {
	// Create a temporary directory for testing
	testDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	// Create some test files
	testFiles := []struct {
		Name string
		Data string
	}{
		{"test1.txt", "test file 1"},
		{"test2.txt", "test file 2"},
		// {"subdir/test3.txt", "test file 3"},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(testDir, tf.Name)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := ioutil.WriteFile(filePath, []byte(tf.Data), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Call the function being tested
	jsonBytes, err := Directory(testDir)
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal the JSON output to a slice of fileInfo structs
	var files []DigestInfo
	if err := json.Unmarshal(jsonBytes, &files); err != nil {
		t.Fatal(err)
	}

	// Assert that the fileInfo slice contains the expected number of elements
	expectedCount := len(testFiles)
	if len(files) != expectedCount {
		t.Fatalf("Expected %d files, got %d", expectedCount, len(files))
	}

	// Assert that each test file is represented in the fileInfo slice
	for _, tf := range testFiles {
		found := false
		for _, f := range files {
			if f.Name == tf.Name {
				found = true
				if f.Sha256 == "" || f.ModTime.IsZero() || f.Mode == 0 {
					t.Errorf("Expected fileInfo to include sha256, mod_time, and mode for %s", f.Name)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected fileInfo for %s, but not found", tf.Name)
		}
	}
}
