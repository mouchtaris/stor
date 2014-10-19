package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/cache"
    "fmt"
    "net"
)

func errorHandler (errors <-chan error) {
    for err := range errors {
        fmt.Printf("error: %s\n", err)
    }
}

func serveNext (s *gocache.Server, l net.Listener, errors chan<- error) {
    conn, err := l.Accept()
    if err != nil {
        errors <- err
        return
    }
    s.GoServe(conn)
}

func serveIncoming (s* gocache.Server, l net.Listener, errors chan<- error) {
    for {
        serveNext(s, l, errors)
    }
}

func newListener (laddr string) net.Listener {
    l, err := net.Listen("tcp", laddr)
    if err != nil {
        panic(err)
    }
    return l
}

func dispatchAll (disp *gocache.Dispatcher, cache *cache.Cache, errors chan<- error) {
    err := disp.DispatchAll(cache)
    if err != nil {
        errors <- err
    }
}

func main () {
    errors := make(chan error, 1)
    cache := cache.NewCache(1)
    dispatcher := gocache.NewDispatcher(1, errors)
    server := gocache.NewServer(20, dispatcher.RequestSink(), errors)
    listener := newListener("0.0.0.0:11000")
    joiner := make(chan uint32, 1)

    go errorHandler(errors)
    go serveIncoming(server, listener, errors)
    go dispatchAll(dispatcher, cache, errors)

    <-joiner
    server.Join()
    server.Close()
    close(errors)
}
