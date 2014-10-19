package command

//
// A common "type" for all commands.
type Command struct {
}

//
// A "delete" command
type Delete struct {
    Command
    key string
}

//
// A "get" command.
type Get struct {
    Command
    Keys []string
}

//
// A "quit" command.
type Quit struct {
    Command
}

//
// A "set" command.
type Set struct {
    Command
    Key, Data string
}

//
// A "stats" command.
type Stats struct {
    Command
}
