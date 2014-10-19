package topcoder_gocache

import (
    "github.com/mouchtaris/topcoder_gocache/command"
)

// A request object, created by the server,
// and handled by the dispatcher.
type Request struct {
    Command command.Command
    Write func (p []byte) (int, error)
    Close func () error
}
