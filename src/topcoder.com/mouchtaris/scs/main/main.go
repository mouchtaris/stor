package main

import (
    "topcoder.com/mouchtaris/scs"
    "topcoder.com/mouchtaris/scs/cache"
    "topcoder.com/mouchtaris/scs/util"
    "fmt"
    "net"
    "os"
    "os/signal"
    "flag"
    "math"
)

const (
    MAX_REQUESTS_QUEUE      = 1
    MAX_SERVER_CONCURRENT   = 1
    MAX_SERVER_QUEUE        = 200
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
    items   uint16
}


func parseCommandLineArguments () Config {
    var items, port uint

    flag.UintVar(&items, "items", 65535, "specify the maximum number of entries allows in the cache")
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

func newRequestsChannel () chan scs.Request {
    return make(chan scs.Request, MAX_REQUESTS_QUEUE)
}

func newServerExecutionManager () scs.ExecutionManager {
    return scs.NewExecutionManager(MAX_SERVER_QUEUE, MAX_SERVER_CONCURRENT)
}

func newServer (requests chan<- scs.Request) *scs.Server {
    eh := util.NewErrorHandler("server")
    return scs.NewServer(requests, eh.ErrorsChannel())
}

func newTCPListener (config Config) net.Listener {
    laddr := fmt.Sprintf("0.0.0.0:%d", config.port)
    l, err := net.Listen("tcp", laddr)
    if err != nil {
        panic(err)
    }
    return l
}

func newConnServer (s *scs.Server, l net.Listener, em scs.ExecutionManager) *scs.ConnServer {
    return scs.NewConnServer(s, l, em)
}

func main () {
    config          := parseCommandLineArguments()
    requests        := newRequestsChannel()
    serverEM        := newServerExecutionManager()
    server          := newServer(requests)
    listener        := newTCPListener(config)
    connserver      := newConnServer(server, listener, serverEM)
    panic("not ready")
    go connserver.ServeAll()
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
    joiner := make(chan uint32, 1)
    stop := make(chan uint32, 1)
    interruption := make(chan os.Signal, 10)
    shutdown := func () {
        stop <- 1
        //server.Join()
        server.Close()
        close(errors)
    }
    defer shutdown()
    signal.Notify(interruption, os.Interrupt)

    go handleInterrupt(interruption, joiner)
    go dispatchAll(dispatcher, cache, errors)

    <-joiner
}
