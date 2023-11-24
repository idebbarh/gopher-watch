package main

import (
	"fmt"

	"github.com/idebbarh/gopher-watch"
)

func main() {
	events := gopherwatch.Watch("/home/ismail/personal/programming/learning/docs/JavaScript/developer.mozilla.org/en-US/docs")

	go func() {
		for {
			select {
			case event := <-events:
				switch true {
				case event.Types.Write:
					file := event.Info.WriteInfo.Name
					fmt.Printf("edited file: %s\n", file)
					// Add your logic for handling write events
				case event.Types.Create:
					file := event.Info.CreateInfo.Name
					fmt.Printf("created file: %s\n", file)
					// Add your logic for handling create events
				case event.Types.Delete:
					file := event.Info.DeleteInfo.Name
					fmt.Printf("deleted file: %s\n", file)
					// Add your logic for handling delete events
				case event.Types.Rename:
					prevName := event.Info.RenameInfo.PrevName
					newName := event.Info.RenameInfo.NewName
					fmt.Printf("file renamed from %s to %s\n", prevName, newName)
					// Add your logic for handling rename events
				}
			}
		}
	}()

	<-make(chan struct{})
}
