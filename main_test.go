package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPasswordFromFile(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expected    string
		expectError bool
	}{
		{
			name:        "simple password",
			fileContent: "mypassword123",
			expected:    "mypassword123",
			expectError: false,
		},
		{
			name:        "password with trailing newline",
			fileContent: "mypassword123\n",
			expected:    "mypassword123",
			expectError: false,
		},
		{
			name:        "password with leading and trailing whitespace",
			fileContent: "  mypassword123  \n",
			expected:    "mypassword123",
			expectError: false,
		},
		{
			name:        "password with multiple trailing newlines",
			fileContent: "mypassword123\n\n\n",
			expected:    "mypassword123",
			expectError: false,
		},
		{
			name:        "empty file",
			fileContent: "",
			expected:    "",
			expectError: false,
		},
		{
			name:        "whitespace only file",
			fileContent: "   \n\t\n  ",
			expected:    "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with the test content
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "password.txt")
			err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0600)
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}

			// Test the function
			result, err := loadPasswordFromFile(tmpFile)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLoadPasswordFromFile_FileNotFound(t *testing.T) {
	_, err := loadPasswordFromFile("/nonexistent/path/to/file.txt")
	if err == nil {
		t.Error("expected error for non-existent file, got none")
	}
}

func TestLoadPasswordFromFile_PermissionDenied(t *testing.T) {
	// Create a temporary file with no read permissions
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "password.txt")
	err := os.WriteFile(tmpFile, []byte("secret"), 0000)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, err = loadPasswordFromFile(tmpFile)
	if err == nil {
		t.Error("expected error for unreadable file, got none")
	}
}
