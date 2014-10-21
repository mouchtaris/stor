package main

import (
    "topcoder.com/mouchtaris/scs"
    "topcoder.com/mouchtaris/scs/cache"
    "fmt"
    "net"
    "os"
    "os/signal"
    "flag"
)

func dispatchAll (disp *scs.Dispatcher, cache *cache.Cache, errors chan<- error) {
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

type Config struct {
    port,
    max_items   uint16


func parseCommandLineArguments () Config {
    items, port uint

    flag.UintVar(&limit, "items", 65535, "specify the maximum number of entries allows in the cache")
    flag.UintVar(&port, "port", 11212, "specify the tcp port to listen to")
    flag.Parse()

    if items > math.MaxUint16 {
        panic(fmt.Sprintf("items value too big: %d", items))
    }
    if port > math.MaxUint16 {
        panic(fmt.Sprintf("port value too big: %d", port))
    }

    return Config { items: uint16(items), port: uint16(port) }
}

func newRequestsChannel () {
    return make(chan Request, MAX_REQUESTS_QUEUE)
}

func newServerExecutionManager () {
    return NewExecutionManager(
        MAX_SERVER_QUEUE,
        MAX_SERVER_CONCURRENT
    )
}

func newServer (requests chan<- Request) {
    eh := NewErrorHandler("server")
    return NewServer(
        requests,
        eh.ErrorsChannel()
    )
}

func newTCPListener (config Config) {
    laddr := fmt.Sprintf("0.0.0.0:%d", config.port)
    l, err := net.Listen("tcp", laddr)
    if err != nil {
        panic(err)
    }
    return l
}

func main () {
    config          := parseCommandLineArguments()
    requests        := newRequestsChannel()
    serverEM        := newServerExecutionManager()
    server          := newServer(requests)
    listener        := newTCPListener(config)
    stop            := make(chan byte, 1)

    go ServeIncomingTCPClients(server, listener, serverEM, stop)
}

func main1 () {
    limit := uint(0)
    port := uint(0)
    flag.UintVar(&limit, "items", 65535, "specify the maximum number of entries allows in the cache")
    flag.UintVar(&port, "port", 11212, "specify the tcp port to listen to")
    flag.Parse()

    errors := make(chan error, 1)
    cache := cache.NewCache(uint32(limit))
    dispatcher := scs.NewDispatcher(1, errors)
    server := scs.NewServer(dispatcher.RequestSink(), errors)
    listener := newListener(fmt.Sprintf("0.0.0.0:%d", port))
    joiner := make(chan uint32, 1)
    stop := make(chan uint32, 1)
    interruption := make(chan os.Signal, 10)
    shutdown := func () {
        stop <- 1
        //server.Join()
        server.Close()
        close(errors)
        listener.Close()
    }
    defer shutdown()
    signal.Notify(interruption, os.Interrupt)

    go handleInterrupt(interruption, joiner)
    go errorHandler(errors)
    go serveIncoming(server, listener, stop, errors)
    go dispatchAll(dispatcher, cache, errors)

    <-joiner
}
