package command

//
// A "set" command.
type Set struct {
    Command
    key, data string
}
