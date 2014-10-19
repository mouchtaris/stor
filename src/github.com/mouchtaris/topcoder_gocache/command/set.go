package command

import (
    "github.com/mouchtaris/topcoder_gocache/cache"
)

//
// A "set" command.
type Set struct {
    Key, Data string
}

//
func (comm *Set) PerformOn (cache *cache.Cache, w WriteBack) error {
    cache.Set(comm.Key, comm.Data)
    w("STORED\r\n")
    return nil
}
