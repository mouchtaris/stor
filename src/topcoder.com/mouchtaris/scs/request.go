package scs

import (
    "topcoder.com/mouchtaris/scs/command"
)

// A request object, created by the server,
// and handled by the dispatcher.
type Request struct {
    Command command.Command
    Write func (p []byte) (int, error)
}
