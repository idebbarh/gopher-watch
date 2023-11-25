package gopherwatch

import (
	"os"
	"testing"
	"time"
)

type Time = time.Time

func TestGetFolderEntriesInfo(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	defer os.RemoveAll(tempDir)

	// Create some test files and directories
	testFiles := []string{"file1.txt", "file2.txt"}
	testDirs := []string{"dir1", "dir2"}

	testModeTime := map[string]Time{}

	for _, file := range testFiles {
		filePath := tempDir + "/" + file
		createdFile, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
		fileInfo, err := createdFile.Stat()
		testModeTime[filePath] = fileInfo.ModTime()
	}

	for _, dir := range testDirs {
		dirPath := tempDir + "/" + dir
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Run the function to get folder entries information
	entriesInfo := make(FolderEntriesInfo)
	GetFolderEntriesInfo(tempDir, entriesInfo)

	// Validate the results

	// Check if the directory itself is present in the entriesInfo
	if _, ok := entriesInfo[tempDir]; !ok {
		t.Errorf("Expected directory %s in entriesInfo, but not found", tempDir)
	}

	// Check if the size of the entires is valid

	if len(entriesInfo[tempDir].Entries) != len(testDirs)+len(testFiles) {
		t.Errorf("Expected %s Entries of size %d but found %d", tempDir, len(testDirs)+len(testFiles), len(entriesInfo[tempDir].Entries))
	}

	// Check if the entries of the directory are present in the entriesInfo
	for _, file := range testFiles {
		filePath := tempDir + "/" + file
		if _, ok := entriesInfo[filePath]; !ok {
			t.Errorf("Expected file %s in entriesInfo, but not found", filePath)
		}
	}

	for _, dir := range testDirs {
		dirPath := tempDir + "/" + dir
		if _, ok := entriesInfo[dirPath]; !ok {
			t.Errorf("Expected directory %s in entriesInfo, but not found", dirPath)
		}
	}

	// Check if the modification time is recorded for files
	for _, file := range testFiles {
		filePath := tempDir + "/" + file
		if !entriesInfo[filePath].ModTime.Equal(testModeTime[filePath]) {
			t.Errorf("Expected modification time to be equal for its created time for file %s", filePath)
		}
	}

	// Check if the isDir field is correctly set
	for _, file := range testFiles {
		filePath := tempDir + "/" + file
		if entriesInfo[filePath].isDir {
			t.Errorf("Expected isDir=false for file %s, but found true", filePath)
		}
	}

	for _, dir := range testDirs {
		dirPath := tempDir + "/" + dir
		if !entriesInfo[dirPath].isDir {
			t.Errorf("Expected isDir=true for directory %s, but found false", dirPath)
		}
	}
}
