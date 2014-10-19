package topcoder_gocache

import (
    "github.com/mouchtaris/topcoder_gocache/lex"
    "github.com/mouchtaris/topcoder_gocache/parser/action"
    "github.com/mouchtaris/topcoder_gocache/parser"
    "github.com/mouchtaris/topcoder_gocache/command"
    "io"
)

type Server struct {
    sem chan uint32
    commands chan<- command.Command
    errors chan<- error
    done chan uint32
    running uint32
}

//
// Construct a new server, which will start at most "backlog"
// concurrent servers.
func NewServer (backlog uint32, commands chan<- command.Command, errors chan<- error) *Server {
    return &Server {
        sem: make(chan uint32, backlog),
        done: make(chan uint32, backlog),
        running: 0,
        commands: commands,
        errors: errors,
    }
}

//
//
func serve (commands chan<- command.Command, input io.ReadCloser) error {
    defer func () {
        input.Close()
    }()

    lexer := lex.NewLexer(input)
    parser := parser.NewParser(lexer)
    parser.RegisterHandler(action.NewSet(commands))
    parser.RegisterHandler(action.NewGet(commands))
    parser.RegisterHandler(action.NewDelete(commands))
    parser.RegisterHandler(action.NewStats(commands))
    parser.RegisterHandler(action.NewQuit(commands))

    err := parser.Parse()
    for ; err == nil; err = parser.Parse() {
    }

    return err
}

//
//
func (server* Server) wrapServing (input io.ReadCloser) {
    err := serve(server.commands, input)
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
func (server *Server) GoServe (input io.ReadCloser) {
    for cont := true; cont; {
        select {
        case server.sem <- 1:
            cont = false
        default:
            server.drainDone()
        }
    }
    server.running++
    go server.wrapServing(input)
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
    close(server.commands)
}
