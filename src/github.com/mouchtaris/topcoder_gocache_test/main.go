package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/cache"
    "os"
    "io"
    "fmt"
)

type NamedWriter func (p []byte) (int, error)

func (w NamedWriter) Write (p []byte) (int, error) {
    return w(p)
}

func (w NamedWriter) Close () error {
    return nil
}

func NewNamedWriter (fname string) NamedWriter {
    return func (p []byte) (int, error) {
        fmt.Printf("[%s]: ", fname)
        bs, err := os.Stdout.Write(p);
        fmt.Println()
        return bs, err
    }
}

type InputOutput struct {
    r io.ReadCloser
    w io.WriteCloser
}

func NewInputOutput (fname string) InputOutput {
    r, err := os.Open(fname)
    if err != nil {
        panic(err)
    }
    return InputOutput {
        r: r,
        w: NewNamedWriter(fname),
    }
}

var Inputs = []InputOutput {
    NewInputOutput("inputs/00.txt"),
    NewInputOutput("inputs/01.txt"),
    NewInputOutput("inputs/02.txt"),
}

func main () {
    errors := make(chan error, 1)
    cache := cache.NewCache()
    disp := gocache.NewDispatcher(1, errors)
    server := gocache.NewServer(20, disp.RequestSink(), errors)

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
        server.GoServe(inp.r, inp.w)
    }

    server.Join()
    server.Close()
    close(errors)
}
