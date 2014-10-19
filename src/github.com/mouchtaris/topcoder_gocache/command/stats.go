package command

import (
    "github.com/mouchtaris/topcoder_gocache/cache"
    "fmt"
)

//
// A "stats" command.
type Stats struct {
}

//
func (comm *Stats) PerformOn (cache *cache.Cache, w WriteBack) error {
    stats := cache.Stats()
    w(fmt.Sprintf(
        "cmd_get %d\r\n" +
        "cmd_set %d\r\n" +
        "get_hits %d\r\n" +
        "get_misses %d\r\n" +
        "delete_hits %d\r\n" +
        "delete_misses %d\r\n" +
        "curr_items %d\r\n" +
        "limit_items %d\r\n" +
        "END\r\n",
        stats.Gets, stats.Sets,
        stats.GetHits, stats.GetMisses,
        stats.DeleteHits, stats.DeleteMisses,
        stats.CurrentItems, stats.Limit))
    return nil
}
