package command

import (
    "github.com/mouchtaris/topcoder_gocache/cache"
    "errors"
)

var ErrQuit = errors.New("quiting")
//
// A "quit" command.
type Quit struct {
}

//
func (comm *Quit) PerformOn (*cache.Cache, WriteBack) error {
    return ErrQuit
}
