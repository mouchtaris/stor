package command

import (
	"topcoder.com/mouchtaris/scs/cache"
)

//
// A "set" command.
type Set struct {
	Key, Data string
}

//
func (comm *Set) PerformOn(cache *cache.Cache, w WriteBack) error {
	err := cache.Set(comm.Key, comm.Data)
	if err != nil {
		return err
	}
	w("STORED\r\n")
	return nil
}
