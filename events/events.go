package gopherwatch

type ChangeType = int

const (
	WRITE ChangeType = iota
	CREATE
	DELETE
	RENAME
	NOCHANGE
)

type EventsInfo struct {
	WriteInfo  struct{ Name string }
	CreateInfo struct{ Name string }
	DeleteInfo struct{ Name string }
	RenameInfo struct {
		PrevName string
		NewName  string
		IsDir    bool
	}
}

type EventsType struct {
	Write  bool
	Create bool
	Delete bool
	Rename bool
}
type Event struct {
	Types EventsType
	Info  *EventsInfo
}
