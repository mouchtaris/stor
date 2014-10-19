package command

import (
    "topcoder.com/mouchtaris/scs/cache"
)

type WriteBack func (string) error
//
// A command is a command which operates on the backing cache.
type Command interface {
    PerformOn (cach *cache.Cache, w WriteBack) error
}
