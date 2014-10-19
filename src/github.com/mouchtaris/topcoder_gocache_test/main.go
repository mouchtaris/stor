package main

import (
    gocache "github.com/mouchtaris/topcoder_gocache"
    "github.com/mouchtaris/topcoder_gocache/cache"
    "fmt"
    "net"
    "os"
    "os/signal"
    "flag"
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

func serveIncoming (s* gocache.Server, l net.Listener, stop <-chan uint32, errors chan<- error) {
    for {
        select {
        case <-stop:
            break;
        default:
            serveNext(s, l, errors)
        }
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

func handleInterrupt (interruption <-chan os.Signal, joiner chan<- uint32) {
    <-interruption
    fmt.Println("Interrupt signal received -- Byebye")
    joiner <- 1
}

func main () {
    limit := uint(0)
    port := uint(0)
    flag.UintVar(&limit, "items", 65535, "specify the maximum number of entries allows in the cache")
    flag.UintVar(&port, "port", 11212, "specify the tcp port to listen to")
    flag.Parse()

    errors := make(chan error, 1)
    cache := cache.NewCache(uint32(limit))
    dispatcher := gocache.NewDispatcher(1, errors)
    server := gocache.NewServer(20, dispatcher.RequestSink(), errors)
    listener := newListener(fmt.Sprintf("0.0.0.0:%d", port))
    joiner := make(chan uint32, 1)
    stop := make(chan uint32, 1)
    interruption := make(chan os.Signal, 10)
    shutdown := func () {
        stop <- 1
        server.Join()
        server.Close()
        close(errors)
    }
    defer shutdown()
    signal.Notify(interruption, os.Interrupt)

    go handleInterrupt(interruption, joiner)
    go errorHandler(errors)
    go serveIncoming(server, listener, stop, errors)
    go dispatchAll(dispatcher, cache, errors)

    <-joiner
}
