package command

import (
    "topcoder.com/mouchtaris/scs/cache"
)

//
// A "get" command.
type Get struct {
    Keys []string
}

//
func (comm *Get) PerformOn (cache *cache.Cache, w WriteBack) error {
    for _, key := range comm.Keys {
        data, ok := cache.Get(key)
        if ok {
            w("VALUE ")
            w(key)
            w("\r\n")
            w(data)
            w("\r\n")
        }
    }
    w("END\r\n")
    return nil
}
