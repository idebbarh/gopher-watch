package gopherwatch

import (
	"fmt"
	"os"
	"time"
)

func listener(watchingPath string, events chan Event) {
	fmt.Printf("listening on : %s\n", watchingPath)
	prevFolderEntriesInfo := FolderEntriesInfo{}
	getFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)

	for {
		curFolderEntriesInfo := FolderEntriesInfo{}
		getFolderEntriesInfo(watchingPath, curFolderEntriesInfo)

		isSomethingChange, changeType, eventInfo := entriesScanner(watchingPath, prevFolderEntriesInfo, curFolderEntriesInfo)

		if isSomethingChange {
			prevFolderEntriesInfo = make(FolderEntriesInfo)
			getFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)
			events <- Event{Types: EventsType{Write: changeType == WRITE, Create: changeType == CREATE, Delete: changeType == DELETE, Rename: changeType == RENAME}, Info: &eventInfo}
		}

		time.Sleep(1 * time.Second)
	}
}

func Watch(watchingPath string) chan Event {
	fi, err := os.Stat(watchingPath)
	if err != nil {
		fmt.Printf("ERROR: Could not get info of %s : %v", watchingPath, err)
		os.Exit(1)
	}

	mode := fi.Mode()

	if !mode.IsDir() {
		fmt.Printf("Error: could not listener to this path because its not a folder")
		os.Exit(1)
	}
	events := make(chan Event)

	go listener(watchingPath, events)

	return events
}
