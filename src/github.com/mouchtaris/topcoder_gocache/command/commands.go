package command

import (
    "github.com/mouchtaris/topcoder_gocache/cache"
    "fmt"
)

var log = func (format string, args... interface{}) {
    fmt.Print("[commands] ")
    fmt.Printf(format, args...)
    fmt.Println()
}

type Command interface {
    PerformOn (*cache.Cache) error
}

//
// A "delete" command
type Delete struct {
    Key string
}

//
func (comm *Delete) PerformOn (cache *cache.Cache) error {
    log("performing delete %s", comm)
    // TODO implement
    return nil
}

//
// A "get" command.
type Get struct {
    Keys []string
}

//
func (comm *Get) PerformOn (cache *cache.Cache) error {
    log("performing get %s", comm)
    // TODO implement
    return nil
}

//
// A "quit" command.
type Quit struct {
}

//
func (comm *Quit) PerformOn (cache *cache.Cache) error {
    log("performing quit %s", comm)
    // TODO implement
    return nil
}

//
// A "set" command.
type Set struct {
    Key, Data string
}

//
func (comm *Set) PerformOn (cache *cache.Cache) error {
    log("performing set %s", comm)
    // TODO implement
    return nil
}

//
// A "stats" command.
type Stats struct {
}

//
func (comm *Stats) PerformOn (cache *cache.Cache) error {
    log("performing stats %s", comm)
    // TODO implement
    return nil
}
