package gopherwatch

import (
	"fmt"
	"time"
)

func Listener(watchingPath string, events chan Event) {
	fmt.Printf("listening on : %s\n", watchingPath)
	prevFolderEntriesInfo := FolderEntriesInfo{}
	GetFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)

	for {
		curFolderEntriesInfo := FolderEntriesInfo{}
		GetFolderEntriesInfo(watchingPath, curFolderEntriesInfo)
		isSomethingChange, changeType, eventInfo := EntriesScanner(watchingPath, prevFolderEntriesInfo, curFolderEntriesInfo)

		if isSomethingChange {
			prevFolderEntriesInfo = make(FolderEntriesInfo)
			GetFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)
			events <- Event{Types: EventsType{Write: changeType == WRITE, Create: changeType == CREATE, Delete: changeType == DELETE, Rename: changeType == RENAME}, Info: &eventInfo}
		}

		time.Sleep(1 * time.Second)
	}
}
