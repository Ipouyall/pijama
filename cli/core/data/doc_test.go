package data

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestNewDocumentExists(t *testing.T) {
	// Prepare a test file
	content := "Test file content"
	tmpFile, err := ioutil.TempFile("", "testfile*.txt")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	err = ioutil.WriteFile(tmpFile.Name(), []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write content to test file: %v", err)
	}

	// Call the NewDocument function
	id := 123
	doc, err := NewDocument(id, tmpFile.Name())

	// Check for errors
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check the document fields
	if doc.ID != id {
		t.Errorf("Unexpected document ID: got %d, expected %d", doc.ID, id)
	}
	if doc.Filename != tmpFile.Name() {
		t.Errorf("Unexpected document filename: got %s, expected %s", doc.Filename, tmpFile.Name())
	}
	if doc.Content != content {
		t.Errorf("Unexpected document content: got %s, expected %s", doc.Content, content)
	}
}

func TestNewDocumentNotExists(t *testing.T) {
	// Call the NewDocument function with a non-existent file
	id := 123
	filename := "nonexistent.txt"
	_, err := NewDocument(id, filename)

	// Check for expected error
	if !os.IsNotExist(err) {
		t.Errorf("Expected error for non-existent file, got: %v", err)
	}
}
