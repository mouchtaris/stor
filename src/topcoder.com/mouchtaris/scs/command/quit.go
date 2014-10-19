package command

import (
	"errors"
	"topcoder.com/mouchtaris/scs/cache"
)

var ErrQuit = errors.New("quiting")

//
// A "quit" command.
type Quit struct {
}

//
func (comm *Quit) PerformOn(*cache.Cache, WriteBack) error {
	return ErrQuit
}
