package command

import (
    "topcoder.com/mouchtaris/scs/cache"
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
