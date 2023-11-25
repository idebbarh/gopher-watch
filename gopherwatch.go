package gopherwatch

import (
	"fmt"
	"os"
)

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

	go Listener(watchingPath, events)

	return events
}
