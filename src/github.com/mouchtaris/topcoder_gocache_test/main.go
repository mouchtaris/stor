package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/cache"
    "os"
    "io"
    "fmt"
)

func open (fname string) io.ReadCloser {
    r, err := os.Open(fname)
    if err != nil {
        return nil
    }
    return r
}

var Inputs = []io.ReadCloser {
    open("inputs/00.txt"),
    open("inputs/01.txt"),
    open("inputs/02.txt"),
}

func main () {
    errors := make(chan error, 1)
    cache := cache.NewCache()
    disp := gocache.NewDispatcher(1, errors)
    server := gocache.NewServer(20, disp.CommandsChannel(), errors)

    go func () {
        for err := range errors {
            fmt.Printf("error: %s\n", err)
        }
    }()

    go func () {
        err := disp.DispatchAll(cache)
        if err != nil {
            errors <- err
        }
    }()

    for _, inp := range Inputs {
        server.GoServe(inp)
    }

    server.Join()
    server.Close()
    close(errors)
}
