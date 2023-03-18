package digester

import (
	"os"
	"testing"
)

func TestFileDigest(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "fileinfo-test")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write some content to the file
	content := "Hello, world!"
	if _, err := file.Write([]byte(content)); err != nil {
		t.Fatalf("Error writing content to temporary file: %v", err)
	}

	// Close the file
	if err := file.Close(); err != nil {
		t.Fatalf("Error closing temporary file: %v", err)
	}

	// stat the file to get the file info
	fileInfo, err := os.Stat(file.Name())
	if err != nil {
		t.Fatalf("Error getting file info: %v", err)
	}

	// Get the file digest info
	digestInfo, err := File(file.Name(), fileInfo)
	if err != nil {
		t.Fatalf("Error getting file info: %v", err)
	}

	// Check the fields of the FileInfo struct
	if digestInfo.Name != file.Name() {
		t.Errorf("Expected name %q, but got %q", file.Name(), digestInfo.Name)
	}

	if digestInfo.ModTime.Unix() != fileInfo.ModTime().Unix() {
		t.Errorf("Expected mod time %v, but got %v", fileInfo.ModTime().Unix(), digestInfo.ModTime)
	}

	if digestInfo.Mode != fileInfo.Mode() {
		t.Errorf("Expected mode %v, but got %v", fileInfo.Mode(), digestInfo.Mode)
	}

	if len(digestInfo.Sha256) != 64 {
		t.Errorf("Expected sha256 digest length 64, but got %v", len(digestInfo.Sha256))
	}

	// Check that the sha256 digest matches the expected value
	expectedDigest := "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"
	if digestInfo.Sha256 != string(expectedDigest) {
		t.Errorf("Expected sha256 digest %v, but got %v", expectedDigest, digestInfo.Sha256)
	}
}
