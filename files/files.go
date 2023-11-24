package gopherwatch

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type DirInfo struct {
	Entries []string
	ModTime time.Time
	isDir   bool
}

type FolderEntriesInfo = map[string]*DirInfo

func getFolderEntriesInfo(curPath string, entriesInfo FolderEntriesInfo) {
	entries, err := os.ReadDir(curPath)
	if err != nil {
		fmt.Printf("Error: could not get the entries of: %s: %s", curPath, err)
		os.Exit(1)
	}

	entriesInfo[curPath] = &DirInfo{}
	entriesInfo[curPath].Entries = []string{}
	entriesInfo[curPath].isDir = true

	for _, e := range entries {
		entryInfo, err := e.Info()
		if err != nil {
			fmt.Printf("ERROR: Could not get info of %s : %v", e.Name(), err)
			os.Exit(1)
		}

		entryPath := curPath + "/" + entryInfo.Name()

		entriesInfo[curPath].Entries = append(entriesInfo[curPath].Entries, entryPath)

		if entryInfo.IsDir() {
			getFolderEntriesInfo(entryPath, entriesInfo)
			continue
		}

		entriesInfo[entryPath] = &DirInfo{}
		entriesInfo[entryPath].ModTime = entryInfo.ModTime()
		entriesInfo[entryPath].isDir = false
	}
}

func entriesScanner(watchingPath string, prevFolderEntriesInfo FolderEntriesInfo, curFolderEntriesInfo FolderEntriesInfo) (bool, ChangeType, EventsInfo) {
	prevWatchingPathInfo, ok := prevFolderEntriesInfo[watchingPath]
	eventInfo := EventsInfo{}

	if !ok {
		for path := range prevFolderEntriesInfo {
			_, ok := curFolderEntriesInfo[path]
			if !ok && len(strings.Split(path, "/")) == len(strings.Split(watchingPath, "/")) {
				eventInfo.RenameInfo.PrevName = path
				eventInfo.RenameInfo.NewName = watchingPath
				eventInfo.RenameInfo.IsDir = true
				break
			}
		}

		return true, RENAME, eventInfo
	}

	curWatchingPathInfo := curFolderEntriesInfo[watchingPath]
	if len(curWatchingPathInfo.Entries) < len(prevWatchingPathInfo.Entries) {
		for path := range prevFolderEntriesInfo {
			_, ok := curFolderEntriesInfo[path]
			if !ok {
				eventInfo.DeleteInfo.Name = path
				break
			}
		}
		return true, DELETE, eventInfo
	}

	if len(curWatchingPathInfo.Entries) > len(prevWatchingPathInfo.Entries) {
		for path := range curFolderEntriesInfo {
			_, ok := prevFolderEntriesInfo[path]
			if !ok {
				eventInfo.CreateInfo.Name = path
				break
			}
		}
		return true, CREATE, eventInfo
	}

	if len(curWatchingPathInfo.Entries) == len(prevWatchingPathInfo.Entries) {
		for _, curEntryPath := range curWatchingPathInfo.Entries {
			curEntryInfo := curFolderEntriesInfo[curEntryPath]
			if curEntryInfo.isDir {
				isSomethingChange, changeType, eventInfo := entriesScanner(curEntryPath, prevFolderEntriesInfo, curFolderEntriesInfo)
				if isSomethingChange {
					return isSomethingChange, changeType, eventInfo
				}
			} else {
				prevEntryInfo, ok := prevFolderEntriesInfo[curEntryPath]
				if !ok || prevEntryInfo.ModTime.Second() != curEntryInfo.ModTime.Second() {
					if ok {
						eventInfo.WriteInfo.Name = curEntryPath
						return true, WRITE, eventInfo
					} else {
						for path := range prevFolderEntriesInfo {
							_, ok := curFolderEntriesInfo[path]
							if !ok {
								eventInfo.RenameInfo.PrevName = path
								eventInfo.RenameInfo.NewName = curEntryPath
								eventInfo.RenameInfo.IsDir = true
								break
							}
						}
						return true, RENAME, eventInfo
					}
				}
			}
		}
	}

	return false, NOCHANGE, eventInfo
}
