package topcoder_gocache

import (
    "github.com/mouchtaris/topcoder_gocache/command"
    "github.com/mouchtaris/topcoder_gocache/cache"
)

type Dispatcher struct {
    commands chan command.Command
    errors chan<- error
}

//
// Construct a new dispatcher, with the given backlog and
// the given error reporting stream.
// The backlog is the ammount of commands that can be
// buffered and wait for execution.
func NewDispatcher (backlog uint32, errors chan<- error) *Dispatcher {
    return &Dispatcher {
        commands: make(chan command.Command, backlog),
        errors: errors,
    }
}

//
// Fetch commands from the CommandsChannel and perform them one
// by one.
func (disp *Dispatcher) DispatchAll (cach *cache.Cache) error {
    for comm := range disp.commands {
        err := comm.PerformOn(cach)
        if err != nil {
            disp.errors <- err
        }
    }
    return nil
}

func (disp *Dispatcher) CommandsChannel () chan<- command.Command {
    return disp.commands
}
