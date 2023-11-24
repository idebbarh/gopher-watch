# GopherWatch

GopherWatch is a Go library for efficient real-time file system monitoring, providing easy tracking of changes in directories and subfolders.

## Installation

```bash
go get https://github.com/idebbarh/gopher-watch
```

## Usage

```go
import (
	"fmt"
    "github.com/idebbarh/gopher-watch"
)

// ...

// Example usage:
events := gopherwatch.Watch(path)

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
```

## License

This library is open-source and provided under the MIT license. Please refer to the [LICENSE](./LICENSE) file for detailed information regarding the terms and conditions of use, modification, and distribution.
