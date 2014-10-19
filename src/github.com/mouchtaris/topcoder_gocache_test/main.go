package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/cache"
    "github.com/mouchtaris/topcoder_gocache/command"
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
    //disp := gocache.NewDispatcher(1, errors)
    commands := make(chan command.Command, 1)
    server := gocache.NewServer(20, commands, errors)

    go func () {
        for comm  := range commands {
            comm.PerformOn(cache)
        }
    }()

    go func () {
        for err := range errors {
            fmt.Printf("error: %s\n", err)
        }
    }()

//    go func () {
//        err := disp.DispatchAll(cache)
//        if err != nil {
//            errors <- err
//        }
//        joiner <- 1
//    }()

    for _, inp := range Inputs {
        server.GoServe(inp)
    }

    server.Join()
    close(commands)
    close(errors)
}
