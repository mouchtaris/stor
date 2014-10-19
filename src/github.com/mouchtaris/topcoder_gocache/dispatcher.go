package topcoder_gocache

import (
    "github.com/mouchtaris/topcoder_gocache/command"
    "github.com/mouchtaris/topcoder_gocache/cache"
)

type Dispatcher func (command.Command, *cache.Cache) error
