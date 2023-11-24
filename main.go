package gopherwatch

import (
	"fmt"
	"os"

	events "github.com/idebbarh/gopher-watch/events"
	watcher "github.com/idebbarh/gopher-watch/watcher"
)

func Watch(watchingPath string) chan events.Event {
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
	events := make(chan events.Event)

	go watcher.Listener(watchingPath, events)

	return events
}
