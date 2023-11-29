package gopherwatch

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

type TestType struct {
	EntryName string
	IsDir     bool
	Action    int
}

type ExpectedType struct {
	ChangeType int
	Name       string
}

type TestEntry struct {
	Test     TestType
	Expected ExpectedType
}

type TestEntries = []TestEntry

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

func TestEntriesScanner(t *testing.T) {
	tmpDir := t.TempDir()
	defer os.RemoveAll(tmpDir)
	testFiles := []string{"file1.txt", "file2.html", "file3.txt"}
	testDirs := []string{"dir1", "dir2", "dir3"}
	testEntries := TestEntries{}

	for _, file := range testFiles {
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: file, IsDir: false, Action: CREATE}, Expected: ExpectedType{ChangeType: CREATE, Name: tmpDir + "/" + file}})
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: file, IsDir: false, Action: WRITE}, Expected: ExpectedType{ChangeType: WRITE, Name: tmpDir + "/" + file}})
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: file, IsDir: false, Action: DELETE}, Expected: ExpectedType{ChangeType: DELETE, Name: tmpDir + "/" + file}})
	}

	for _, dir := range testDirs {
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: dir, IsDir: true, Action: CREATE}, Expected: ExpectedType{ChangeType: CREATE, Name: tmpDir + "/" + dir}})
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: dir, IsDir: true, Action: DELETE}, Expected: ExpectedType{ChangeType: DELETE, Name: tmpDir + "/" + dir}})
	}

	prevFolderEntriesInfo := make(FolderEntriesInfo)
	GetFolderEntriesInfo(tmpDir, prevFolderEntriesInfo)
	curFolderEntriesInfo := make(FolderEntriesInfo)
	GetFolderEntriesInfo(tmpDir, curFolderEntriesInfo)
	isSomethingChanged, _, _ := EntriesScanner(tmpDir, prevFolderEntriesInfo, curFolderEntriesInfo)

	if isSomethingChanged {
		t.Error("Expected isSomethingChanged=false, but found true")
	}

	numToChangeType := map[int]string{0: "WRITE", 1: "CREATE", 2: "DELETE", 3: "RENAME", 4: "NOCHANGE"}

	// test create event
	for _, entry := range testEntries {
		entryPath := tmpDir + "/" + entry.Test.EntryName

		if entry.Test.Action == CREATE {
			if entry.Test.IsDir {
				err := os.Mkdir(entryPath, os.ModePerm)
				if err != nil {
					t.Fatalf("Failed to create test dir %s: %v", entry.Test.EntryName, err)
				}
			} else {
				_, err := os.Create(entryPath)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", entry.Test.EntryName, err)
				}
				time.Sleep(time.Second)
			}
		} else if entry.Test.Action == WRITE {
			err := os.WriteFile(entryPath, []byte("new line"), os.ModePerm)
			if err != nil {
				t.Fatalf("Failed to write to file %s: %v", entry.Test.EntryName, err)
			}
		} else if entry.Test.Action == DELETE {
			err := os.Remove(entryPath)
			if err != nil {
				t.Fatalf("Failed to delete test entry %s: %v", entry.Test.EntryName, err)
			}
		}

		curFolderEntriesInfo = make(FolderEntriesInfo)
		GetFolderEntriesInfo(tmpDir, curFolderEntriesInfo)

		isSomethingChanged, changeType, eventInfo := EntriesScanner(tmpDir, prevFolderEntriesInfo, curFolderEntriesInfo)

		if !isSomethingChanged {
			t.Error("Expected isSomethingChanged=true, but found false")
		}

		if changeType != entry.Expected.ChangeType {
			t.Errorf("Expected changeType=%s, but found %s", numToChangeType[entry.Expected.ChangeType], numToChangeType[changeType])
		}

		if entry.Test.Action == CREATE && eventInfo.CreateInfo.Name != entry.Expected.Name {
			t.Errorf("Expected eventInfo.CreateInfo.name=%s, but found %s", entry.Expected.Name, eventInfo.CreateInfo.Name)
		} else if entry.Test.Action == DELETE && eventInfo.DeleteInfo.Name != entry.Expected.Name {
			t.Errorf("Expected eventInfo.DeleteInfo.name=%s, but found %s", entry.Expected.Name, eventInfo.WriteInfo.Name)
		} else if entry.Test.Action == WRITE && eventInfo.WriteInfo.Name != entry.Expected.Name {
			t.Errorf("Expected eventInfo.WriteInfo.name=%s, but found %s", entry.Expected.Name, eventInfo.WriteInfo.Name)
		}

		prevFolderEntriesInfo = make(FolderEntriesInfo)
		GetFolderEntriesInfo(tmpDir, prevFolderEntriesInfo)
	}

	for _, file := range testFiles {
		filePath := tmpDir + "/" + file
		_, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	for _, dir := range testDirs {
		dirPath := tmpDir + "/" + dir
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	prevFolderEntriesInfo = make(FolderEntriesInfo)
	GetFolderEntriesInfo(tmpDir, prevFolderEntriesInfo)

	for i, entry := range append(testFiles, testDirs...) {
		entryPath := tmpDir + "/" + entry
		newEntryPath := tmpDir + "/" + entry + fmt.Sprintf("%d", i)
		err := os.Rename(entryPath, newEntryPath)
		isDirEntry := len(strings.Split(entry, ".")) == 1
		if err != nil {
			t.Fatalf("Failed to rename entry %s to %s: %v", entryPath, newEntryPath, err)
		}

		curFolderEntriesInfo = make(FolderEntriesInfo)
		GetFolderEntriesInfo(tmpDir, curFolderEntriesInfo)
		isSomethingChanged, changeType, eventInfo := EntriesScanner(tmpDir, prevFolderEntriesInfo, curFolderEntriesInfo)

		if !isSomethingChanged {
			t.Error("Expected isSomethingChanged=true, but found false")
		}

		if changeType != RENAME {
			t.Errorf("Expected changeType=%s, but found %s", numToChangeType[RENAME], numToChangeType[changeType])
		}

		if eventInfo.RenameInfo.IsDir != isDirEntry {
			t.Errorf("Expected eventInfo.RenameInfo.IsDir=%v, but found %v", isDirEntry, eventInfo.RenameInfo.IsDir)
		}

		if eventInfo.RenameInfo.PrevName != entryPath {
			t.Errorf("Expected eventInfo.RenameInfo.PrevName=%s, but found %s", entryPath, eventInfo.RenameInfo.PrevName)
		}

		if eventInfo.RenameInfo.NewName != newEntryPath {
			t.Errorf("Expected eventInfo.RenameInfo.PrevName=%s, but found %s", newEntryPath, eventInfo.RenameInfo.NewName)
		}

		prevFolderEntriesInfo = make(FolderEntriesInfo)
		GetFolderEntriesInfo(tmpDir, prevFolderEntriesInfo)
	}
}

func TestListener(t *testing.T) {
	tmpDir := t.TempDir()
	var curTest TestEntry

	testEvents := make(chan Event)

	go Listener(tmpDir, testEvents)

	go func() {
		for {
			select {
			case event := <-testEvents:
				switch curTest.Expected.ChangeType {
				case CREATE:
					if !event.Types.Create {
						t.Error("Expected event.Types.Create=true , but found false")
					}
					if event.Info.CreateInfo.Name != curTest.Expected.Name {
						t.Errorf("Expected event.Info.CreateInfo.Name=%s , but found %s", curTest.Expected.Name, event.Info.CreateInfo.Name)
					}

				case WRITE:
					if !event.Types.Write {
						t.Error("Expected event.Types.Write=true , but found false")
					}

				case DELETE:
					if !event.Types.Delete {
						t.Error("Expected event.Types.Delete=true , but found false")
					}
				}
			}
		}
	}()

	testFiles := []string{"file1.txt", "file2.html", "file3.txt"}
	testDirs := []string{"dir1", "dir2", "dir3"}

	testEntries := TestEntries{}

	for _, file := range testFiles {
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: file, IsDir: false, Action: CREATE}, Expected: ExpectedType{ChangeType: CREATE, Name: tmpDir + "/" + file}})
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: file, IsDir: false, Action: WRITE}, Expected: ExpectedType{ChangeType: WRITE, Name: tmpDir + "/" + file}})
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: file, IsDir: false, Action: DELETE}, Expected: ExpectedType{ChangeType: DELETE, Name: tmpDir + "/" + file}})
	}

	for _, dir := range testDirs {
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: dir, IsDir: true, Action: CREATE}, Expected: ExpectedType{ChangeType: CREATE, Name: tmpDir + "/" + dir}})
		testEntries = append(testEntries, TestEntry{Test: TestType{EntryName: dir, IsDir: true, Action: DELETE}, Expected: ExpectedType{ChangeType: DELETE, Name: tmpDir + "/" + dir}})
	}

	for _, testEntry := range testEntries {
		entryPath := tmpDir + "/" + testEntry.Test.EntryName
		curTest = testEntry
		if testEntry.Test.Action == CREATE {
			if testEntry.Test.IsDir {
				err := os.Mkdir(entryPath, os.ModePerm)
				if err != nil {
					t.Fatalf("Failed to create test dir %s: %v", testEntry.Test.EntryName, err)
				}
			} else {
				_, err := os.Create(entryPath)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", testEntry.Test.EntryName, err)
				}
			}
		} else if testEntry.Test.Action == WRITE {
			err := os.WriteFile(entryPath, []byte("new line"), os.ModePerm)
			if err != nil {
				t.Fatalf("Failed to write to file %s: %v", testEntry.Test.EntryName, err)
			}
		} else if testEntry.Test.Action == DELETE {
			err := os.Remove(entryPath)
			if err != nil {
				t.Fatalf("Failed to delete test entry %s: %v", testEntry.Test.EntryName, err)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
