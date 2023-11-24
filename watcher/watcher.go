package gopherwatch

import (
	"fmt"
	"time"

	e "github.com/idebbarh/gopher-watch/events"
	f "github.com/idebbarh/gopher-watch/files"
)

func Listener(watchingPath string, events chan e.Event) {
	fmt.Printf("listening on : %s\n", watchingPath)
	prevFolderEntriesInfo := f.FolderEntriesInfo{}
	f.GetFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)

	for {
		curFolderEntriesInfo := f.FolderEntriesInfo{}
		f.GetFolderEntriesInfo(watchingPath, curFolderEntriesInfo)
		isSomethingChange, changeType, eventInfo := f.EntriesScanner(watchingPath, prevFolderEntriesInfo, curFolderEntriesInfo)

		if isSomethingChange {
			prevFolderEntriesInfo = make(f.FolderEntriesInfo)
			f.GetFolderEntriesInfo(watchingPath, prevFolderEntriesInfo)
			events <- e.Event{Types: e.EventsType{Write: changeType == e.WRITE, Create: changeType == e.CREATE, Delete: changeType == e.DELETE, Rename: changeType == e.RENAME}, Info: &eventInfo}
		}

		time.Sleep(1 * time.Second)
	}
}
