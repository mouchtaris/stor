package topcoder_gocache

import (
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/parser/action"
    "github.com/mouchtaris/topcoder_gocache/parser"
    "io"
)

type Server struct {
    sem chan uint32
    requests chan<- Request
    errors chan<- error
    done chan uint32
    running uint32
}

//
// Construct a new server, which will start at most "backlog"
// concurrent servers.
func NewServer (backlog uint32, requests chan<- Request, errors chan<- error) *Server {
    return &Server {
        sem: make(chan uint32, backlog),
        done: make(chan uint32, backlog),
        running: 0,
        requests: requests,
        errors: errors,
    }
}

//
//
func serve (requests chan<- Request, input io.ReadCloser, output io.WriteCloser) error {
    defer input.Close()

    lexer := lex.NewLexer(input)
    parser := parser.NewParser(lexer)
    parser.RegisterHandler(action.Set    { })
    parser.RegisterHandler(action.Get    { })
    parser.RegisterHandler(action.Delete { })
    parser.RegisterHandler(action.Stats  { })
    parser.RegisterHandler(action.Quit   { })

    comm, err := parser.Parse()
    for ; err == nil; comm, err = parser.Parse() {
        requests <- Request { comm, output.Write, output.Close, }
    }

    return err
}

//
//
func (server* Server) wrapServing (input io.ReadCloser, output io.WriteCloser) {
    err := serve(server.requests, input, output)
    if err != nil {
        server.errors<- err
    }
    server.done <- 1
    <-server.sem
}

//
//
func (server *Server) drainDone () {
    for cont := true; cont; {
        select {
        case <-server.done:
            server.running--
        default:
            cont = false
        }
    }
}

//
// Server (asynchronously) parsing input from the given reader.
// Parsed commands are send to the server's command stream
// for further processing by some other consumer.
func (server *Server) GoServe (input io.ReadCloser, output io.WriteCloser) {
    for cont := true; cont; {
        select {
        case server.sem <- 1:
            cont = false
        default:
            server.drainDone()
        }
    }
    server.running++
    go server.wrapServing(input, output)
}

//
// Wait for all running jobs to finish.
func (server *Server) Join () {
    for i := uint32(0); i < server.running; i++ {
        <-server.done
    }
    server.running = 0
}

//
// Close this server. This closes the underlying
// commands channel and thus serving anything else
// will result in panic.
func (server *Server) Close () {
    close(server.requests)
}
