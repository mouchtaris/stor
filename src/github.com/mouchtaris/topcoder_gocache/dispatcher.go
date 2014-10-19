package topcoder_gocache

import (
    "github.com/mouchtaris/topcoder_gocache/cache"
)

type Dispatcher struct {
    requests chan Request
    errors chan<- error
}

//
// Construct a new dispatcher, with the given backlog and
// the given error reporting stream.
// The backlog is the ammount of commands that can be
// buffered and wait for execution.
func NewDispatcher (backlog uint32, errors chan<- error) *Dispatcher {
    return &Dispatcher {
        requests: make(chan Request, backlog),
        errors: errors,
    }
}

//
// Fetch commands from the CommandsChannel and perform them one
// by one.
func (disp *Dispatcher) DispatchAll (cach *cache.Cache) error {
    for req := range disp.requests {
        writeBack := func (s string) error {
            _, err := req.Write([]byte(s))
            return err
        }
        err := req.Command.PerformOn(cach, writeBack)

        if err != nil {
            disp.errors <- err
            req.Close()
        }
    }
    return nil
}

func (disp *Dispatcher) RequestSink () chan<- Request {
    return disp.requests
}
