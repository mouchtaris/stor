package command

//
// A "delete" command
type Delete struct {
    Key string
}

//
// A "get" command.
type Get struct {
    Keys []string
}

//
// A "quit" command.
type Quit struct {
}

//
// A "set" command.
type Set struct {
    Key, Data string
}

//
// A "stats" command.
type Stats struct {
}
