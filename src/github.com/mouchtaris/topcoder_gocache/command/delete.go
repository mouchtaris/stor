package command

import (
    "github.com/mouchtaris/topcoder_gocache/cache"
)

//
// A "delete" command
type Delete struct {
    Key string
}

//
func (comm *Delete) PerformOn (cache *cache.Cache, w WriteBack) error {
    _, deleted := cache.Delete(comm.Key)
    if deleted {
        w("DELETED\r\n")
    } else {
        w("NOT_FOUND\r\n")
    }
    return nil
}
